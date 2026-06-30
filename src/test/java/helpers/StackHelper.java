package helpers;

import io.qameta.allure.Step;

import java.io.IOException;
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.nio.file.Path;
import java.time.Duration;

public final class StackHelper {

    private static final HttpClient HTTP = HttpClient.newBuilder()
            .connectTimeout(Duration.ofSeconds(5))
            .build();

    private static final Path DEV_ROOT = Path.of(System.getProperty("user.dir")).resolve("..").resolve("dev").normalize();

    private StackHelper() {
    }

    @Step("Stop hub on :4444")
    public static void killHub() {
        execShell("lsof -nP -iTCP:4444 -sTCP:LISTEN -t | xargs kill 2>/dev/null || true");
    }

    @Step("Start hub detached")
    public static void startHubDetached() {
        execDetached(DEV_ROOT.resolve("scripts/start-selenoid.sh").toString());
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

    private static void execDetached(String scriptPath) {
        try {
            new ProcessBuilder("bash", scriptPath)
                    .directory(DEV_ROOT.toFile())
                    .redirectErrorStream(true)
                    .start();
        } catch (IOException e) {
            throw new IllegalStateException("Failed to start: " + scriptPath, e);
        }
    }
}
