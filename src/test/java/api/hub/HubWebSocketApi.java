package api.hub;

import io.qameta.allure.Step;

import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.WebSocket;
import java.time.Duration;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.TimeUnit;

public final class HubWebSocketApi {

    private HubWebSocketApi() {
    }

    @Step("WebSocket {path} — expect open")
    public static void openBriefly(String path, Duration timeout) throws Exception {
        var opened = new CompletableFuture<Void>();
        var failed = new CompletableFuture<Void>();
        HttpClient.newHttpClient()
                .newWebSocketBuilder()
                .connectTimeout(timeout)
                .buildAsync(hubWebSocketUri(path), new WebSocket.Listener() {
                    @Override
                    public void onOpen(WebSocket webSocket) {
                        opened.complete(null);
                        webSocket.sendClose(WebSocket.NORMAL_CLOSURE, "api-test");
                    }

                    @Override
                    public void onError(WebSocket webSocket, Throwable error) {
                        failed.completeExceptionally(error);
                    }
                })
                .get(timeout.toMillis(), TimeUnit.MILLISECONDS);
        try {
            opened.get(timeout.toMillis(), TimeUnit.MILLISECONDS);
        } catch (Exception e) {
            if (failed.isCompletedExceptionally()) {
                failed.get(timeout.toMillis(), TimeUnit.MILLISECONDS);
            }
            throw e;
        }
    }

    public static URI hubWebSocketUri(String path) {
        var http = HubRequest.baseUri();
        var base = URI.create(http);
        var wsScheme = "https".equalsIgnoreCase(base.getScheme()) ? "wss" : "ws";
        var port = base.getPort();
        var authority = port > 0 ? base.getHost() + ":" + port : base.getHost();
        var normalizedPath = path.startsWith("/") ? path : "/" + path;
        return URI.create(wsScheme + "://" + authority + normalizedPath);
    }
}
