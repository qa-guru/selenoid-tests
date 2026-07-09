package helpers;

import io.qameta.allure.Step;

import java.io.IOException;
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.nio.file.Files;
import java.nio.file.Path;
import java.time.Duration;
import java.util.ArrayList;
import java.util.List;

public final class StackHelper {

    private static final HttpClient HTTP = HttpClient.newBuilder()
            .connectTimeout(Duration.ofSeconds(5))
            .build();

    private static final Path PROJECT_ROOT = Path.of(System.getProperty("user.dir")).toAbsolutePath().normalize();
    private static final Path DEV_ROOT = PROJECT_ROOT.resolve("..").resolve("dev").normalize();
    private static final Path CI_BIN = PROJECT_ROOT.resolve("build").resolve("ci-bin").resolve("selenoid");
    private static final Path CI_BROWSERS = PROJECT_ROOT.resolve("fixtures").resolve("ci-browsers.json");
    private static final Path DEV_START = DEV_ROOT.resolve("scripts").resolve("start-selenoid.sh");

    private StackHelper() {
    }

    @Step("Stop hub on :4444")
    public static void killHub() {
        execShell("lsof -nP -iTCP:4444 -sTCP:LISTEN -t | xargs kill 2>/dev/null || true");
    }

    @Step("Start hub detached")
    public static void startHubDetached() {
        var ciBin = resolveCiHubBinary();
        if (ciBin != null) {
            execDetached(List.of(
                    ciBin.toString(),
                    "-conf", CI_BROWSERS.toString(),
                    "-limit", "3",
                    "-video-recorder-image", "qaguru/video-recorder:latest"));
            return;
        }
        if (!Files.isRegularFile(DEV_START)) {
            throw new IllegalStateException(
                    "No hub binary: missing " + CI_BIN + " and " + DEV_START);
        }
        execDetached(List.of("bash", DEV_START.toString()));
    }

    @Step("Wait until hub /status is ready")
    public static void waitForHubReady(long timeoutMs) throws InterruptedException {
        var deadline = System.currentTimeMillis() + timeoutMs;
        while (System.currentTimeMillis() < deadline) {
            if (hubResponds()) {
                return;
            }
            Thread.sleep(500);
        }
        throw new IllegalStateException("Hub did not become ready on :4444 within " + timeoutMs + "ms");
    }

    @Step("Wait until hub /status is down")
    public static void waitForHubDown(long timeoutMs) throws InterruptedException {
        var deadline = System.currentTimeMillis() + timeoutMs;
        while (System.currentTimeMillis() < deadline) {
            if (!hubResponds()) {
                return;
            }
            Thread.sleep(500);
        }
        throw new IllegalStateException("Hub is still responding on :4444 after kill");
    }

    private static Path resolveCiHubBinary() {
        var fromEnv = System.getenv("SELENOID_BIN");
        if (fromEnv != null && !fromEnv.isBlank()) {
            var path = Path.of(fromEnv);
            if (Files.isExecutable(path)) {
                return path;
            }
        }
        if (Files.isExecutable(CI_BIN) && Files.isRegularFile(CI_BROWSERS)) {
            return CI_BIN;
        }
        return null;
    }

    private static boolean hubResponds() {
        try {
            var request = HttpRequest.newBuilder(URI.create("http://127.0.0.1:4444/status"))
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

    private static void execShell(String command) {
        try {
            var process = new ProcessBuilder("bash", "-lc", command)
                    .redirectErrorStream(true)
                    .start();
            process.waitFor();
        } catch (IOException | InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new IllegalStateException("Shell command failed: " + command, e);
        }
    }

    private static void execDetached(List<String> command) {
        try {
            var builder = new ProcessBuilder(new ArrayList<>(command));
            builder.directory(PROJECT_ROOT.toFile());
            builder.redirectErrorStream(true);
            builder.redirectOutput(ProcessBuilder.Redirect.DISCARD);
            builder.start();
        } catch (IOException e) {
            throw new IllegalStateException("Failed to start: " + command, e);
        }
    }
}
