package helpers;

import config.ConfigReader;
import config.TestConfig;
import io.qameta.allure.Step;

import java.io.IOException;
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.attribute.PosixFilePermission;
import java.time.Duration;
import java.util.ArrayList;
import java.util.EnumSet;
import java.util.List;
import java.util.concurrent.TimeUnit;

public final class CmInstallerHelper {

    public record CmRunResult(int exitCode, String output) {
        public void requireSuccess(String action) {
            if (exitCode != 0) {
                throw new IllegalStateException(
                        action + " failed (exit " + exitCode + "):\n" + output);
            }
        }
    }

    private static final HttpClient HTTP = HttpClient.newBuilder()
            .connectTimeout(Duration.ofSeconds(5))
            .build();

    private static final Path PROJECT_ROOT = Path.of(System.getProperty("user.dir"));
    private static final Path WORKSPACE_ROOT = PROJECT_ROOT.resolve("..").normalize();

    private final TestConfig config;
    private final Path configDir;
    private Path cmBinary;
    private Path browsersJson;

    private CmInstallerHelper(TestConfig config, Path configDir) {
        this.config = config;
        this.configDir = configDir;
    }

    private Path cmBinary() {
        if (cmBinary == null) {
            cmBinary = resolveExecutable(config.cmBinaryPath(), "cm binary");
        }
        return cmBinary;
    }

    private Path browsersJson() {
        if (browsersJson == null) {
            browsersJson = resolveExisting(config.cmBrowsersJson(), "browsers.json");
        }
        return browsersJson;
    }

    public static CmInstallerHelper withTempConfigDir() throws IOException {
        var dir = Files.createTempDirectory("cm-installer-");
        return new CmInstallerHelper(ConfigReader.testConfig, dir);
    }

    public Path configDir() {
        return configDir;
    }

    public Path browsersJsonPath() {
        return configDir.resolve("browsers.json");
    }

    @Step("Stop CM-managed hub and UI")
    public void stopAll() {
        stopUiQuietly();
        stopHubQuietly();
    }

    @Step("cm selenoid configure")
    public CmRunResult configure() {
        return runSelenoid(
                "configure",
                "-c", configDir.toString(),
                "-p", Integer.toString(config.cmHubPort()),
                "-n",
                "-j", browsersJson().toString());
    }

    @Step("cm selenoid start")
    public CmRunResult startHub() throws IOException {
        ensureLinuxBinary("selenoid", config.cmSelenoidBinary(), selenoidRepoDir());
        return runSelenoid(
                "start",
                "-c", configDir.toString(),
                "-p", Integer.toString(config.cmHubPort()),
                "-n",
                "-j", browsersJson().toString(),
                "--selenoid-binary", configDir.resolve("bin/selenoid").toAbsolutePath().toString(),
                "--selenoid-ui-binary", configDir.resolve("bin/selenoid-ui").toAbsolutePath().toString());
    }

    @Step("cm selenoid-ui start")
    public CmRunResult startUi() throws IOException {
        ensureLinuxBinary("selenoid-ui", config.cmSelenoidUiBinary(), selenoidUiRepoDir());
        return runSelenoidUi(
                "start",
                "-c", configDir.toString(),
                "-p", Integer.toString(config.cmUiPort()));
    }

    @Step("cm selenoid status")
    public CmRunResult statusHub() {
        return runSelenoid("status", "-c", configDir.toString(), "-p", Integer.toString(config.cmHubPort()));
    }

    @Step("cm selenoid-ui status")
    public CmRunResult statusUi() {
        return runSelenoidUi("status", "-c", configDir.toString(), "-p", Integer.toString(config.cmUiPort()));
    }

    @Step("Wait until hub /status is ready")
    public void waitForHubReady(long timeoutMs) throws InterruptedException {
        waitForStatusReady(ConfigReader.resolveCmHubUrl(config), timeoutMs, "hub");
    }

    @Step("Wait until UI /status is ready")
    public void waitForUiReady(long timeoutMs) throws InterruptedException {
        waitForStatusReady(ConfigReader.resolveCmUiUrl(config), timeoutMs, "UI");
    }

    public void deleteConfigDir() {
        try {
            if (Files.exists(configDir)) {
                Files.walk(configDir)
                        .sorted((a, b) -> b.compareTo(a))
                        .forEach(path -> {
                            try {
                                Files.deleteIfExists(path);
                            } catch (IOException ignored) {
                                // best-effort cleanup
                            }
                        });
            }
        } catch (IOException e) {
            throw new IllegalStateException("Failed to delete config dir: " + configDir, e);
        }
    }

    @Step("Ensure Linux {binaryName} binary for Docker wrapper")
    private void ensureLinuxBinary(String binaryName, String localBinaryProperty, Path sourceRepo)
            throws IOException {
        var target = configDir.resolve("bin").resolve(binaryName);
        if (Files.isExecutable(target)) {
            return;
        }
        Files.createDirectories(target.getParent());

        if (config.cmUseLocalBinaries() && isLinuxHost()) {
            var local = resolveExecutable(localBinaryProperty, binaryName + " binary");
            Files.copy(local, target);
            makeExecutable(target);
            return;
        }

        crossCompileLinuxBinary(sourceRepo, target);
    }

    private void crossCompileLinuxBinary(Path sourceRepo, Path target) {
        if (!Files.isDirectory(sourceRepo)) {
            throw new IllegalStateException("Source repo not found: " + sourceRepo);
        }

        try {
            var process = new ProcessBuilder("go", "build", "-o", target.toString(), ".")
                    .directory(sourceRepo.toFile())
                    .redirectErrorStream(true);
            var env = process.environment();
            env.put("GOOS", "linux");
            env.put("GOARCH", linuxGoArch());
            env.put("CGO_ENABLED", "0");

            var started = process.start();
            var output = new String(started.getInputStream().readAllBytes(), StandardCharsets.UTF_8);
            if (!started.waitFor(3, TimeUnit.MINUTES)) {
                started.destroyForcibly();
                throw new IllegalStateException("go build timed out in " + sourceRepo);
            }
            if (started.exitValue() != 0) {
                throw new IllegalStateException(
                        "go build failed in " + sourceRepo + " (exit " + started.exitValue() + "):\n" + output);
            }
            makeExecutable(target);
        } catch (IOException | InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new IllegalStateException("Failed to cross-compile in " + sourceRepo + ": " + e.getMessage(), e);
        }
    }

    private static Path selenoidRepoDir() {
        return WORKSPACE_ROOT.resolve("selenoid");
    }

    private static Path selenoidUiRepoDir() {
        return WORKSPACE_ROOT.resolve("selenoid-ui");
    }

    private static boolean isLinuxHost() {
        return System.getProperty("os.name", "").toLowerCase().contains("linux");
    }

    private static String linuxGoArch() {
        var arch = System.getProperty("os.arch", "").toLowerCase();
        return arch.contains("arm") || arch.contains("aarch64") ? "arm64" : "amd64";
    }

    private static void makeExecutable(Path binary) {
        var perms = EnumSet.of(
                PosixFilePermission.OWNER_READ,
                PosixFilePermission.OWNER_WRITE,
                PosixFilePermission.OWNER_EXECUTE);
        try {
            Files.setPosixFilePermissions(binary, perms);
        } catch (UnsupportedOperationException ignored) {
            // Windows
        } catch (IOException e) {
            throw new IllegalStateException("Failed to chmod " + binary, e);
        }
    }

    private void stopHubQuietly() {
        try {
            runSelenoid("stop", "-c", configDir.toString(), "-p", Integer.toString(config.cmHubPort()));
        } catch (RuntimeException ignored) {
            // container may already be stopped
        }
    }

    private void stopUiQuietly() {
        try {
            runSelenoidUi("stop", "-c", configDir.toString(), "-p", Integer.toString(config.cmUiPort()));
        } catch (RuntimeException ignored) {
            // container may already be stopped
        }
    }

    private CmRunResult runSelenoid(String subcommand, String... args) {
        return run("selenoid", subcommand, args);
    }

    private CmRunResult runSelenoidUi(String subcommand, String... args) {
        return run("selenoid-ui", subcommand, args);
    }

    private CmRunResult run(String scope, String subcommand, String... args) {
        var command = new ArrayList<String>();
        command.add(cmBinary().toString());
        command.add(scope);
        command.add(subcommand);
        command.addAll(List.of(args));

        try {
            var process = new ProcessBuilder(command)
                    .redirectErrorStream(true)
                    .start();
            var output = new String(process.getInputStream().readAllBytes(), StandardCharsets.UTF_8);
            if (!process.waitFor(5, TimeUnit.MINUTES)) {
                process.destroyForcibly();
                throw new IllegalStateException(
                        "cm " + scope + " " + subcommand + " timed out after 5 minutes");
            }
            return new CmRunResult(process.exitValue(), output);
        } catch (IOException | InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new IllegalStateException(
                    "Failed to run cm " + scope + " " + subcommand + ": " + e.getMessage(), e);
        }
    }

    private static void waitForStatusReady(String baseUrl, long timeoutMs, String label)
            throws InterruptedException {
        var deadline = System.currentTimeMillis() + timeoutMs;
        var statusUri = URI.create(baseUrl + "status");
        while (System.currentTimeMillis() < deadline) {
            if (statusResponds(statusUri)) {
                return;
            }
            Thread.sleep(500);
        }
        throw new IllegalStateException(label + " did not become ready at " + statusUri
                + " within " + timeoutMs + "ms");
    }

    private static boolean statusResponds(URI statusUri) {
        try {
            var request = HttpRequest.newBuilder(statusUri)
                    .timeout(Duration.ofSeconds(3))
                    .GET()
                    .build();
            var response = HTTP.send(request, HttpResponse.BodyHandlers.discarding());
            return response.statusCode() == 200;
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            return false;
        } catch (IOException e) {
            return false;
        }
    }

    private static Path resolveExecutable(String path, String label) {
        var resolved = resolvePath(path);
        if (!Files.isExecutable(resolved)) {
            throw new IllegalStateException(label + " not found or not executable: " + resolved);
        }
        return resolved;
    }

    private static Path resolveExisting(String path, String label) {
        var resolved = resolvePath(path);
        if (!Files.isRegularFile(resolved)) {
            throw new IllegalStateException(label + " not found: " + resolved);
        }
        return resolved;
    }

    private static Path resolvePath(String path) {
        var candidate = Path.of(path);
        return candidate.isAbsolute() ? candidate : PROJECT_ROOT.resolve(candidate).normalize();
    }
}
