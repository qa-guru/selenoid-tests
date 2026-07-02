package helpers;

import config.ConfigReader;
import config.TestConfig;
import io.qameta.allure.Step;

import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.TimeUnit;

public final class CmCli {

    private static final Path PROJECT_ROOT = Path.of(System.getProperty("user.dir"));

    private CmCli() {
    }

    @Step("cm version")
    public static CmInstallerHelper.CmRunResult version() {
        return run(List.of("version"));
    }

    @Step("cm --help")
    public static CmInstallerHelper.CmRunResult help() {
        return run(List.of("--help"));
    }

    private static CmInstallerHelper.CmRunResult run(List<String> args) {
        var config = ConfigReader.testConfig;
        var binary = resolveBinary(config);
        var command = new ArrayList<String>();
        command.add(binary.toString());
        command.addAll(args);

        try {
            var process = new ProcessBuilder(command)
                    .redirectErrorStream(true)
                    .start();
            var output = new String(process.getInputStream().readAllBytes(), StandardCharsets.UTF_8);
            if (!process.waitFor(2, TimeUnit.MINUTES)) {
                process.destroyForcibly();
                throw new IllegalStateException("cm " + String.join(" ", args) + " timed out");
            }
            return new CmInstallerHelper.CmRunResult(process.exitValue(), output);
        } catch (IOException | InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new IllegalStateException("Failed to run cm " + String.join(" ", args) + ": " + e.getMessage(), e);
        }
    }

    private static Path resolveBinary(TestConfig config) {
        var candidate = Path.of(config.cmBinaryPath());
        var resolved = candidate.isAbsolute() ? candidate : PROJECT_ROOT.resolve(candidate).normalize();
        if (!Files.isExecutable(resolved)) {
            throw new IllegalStateException("cm binary not found or not executable: " + resolved);
        }
        return resolved;
    }
}
