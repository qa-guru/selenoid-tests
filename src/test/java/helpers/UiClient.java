package helpers;

import com.google.gson.JsonElement;
import com.google.gson.JsonObject;
import com.google.gson.JsonParser;
import io.qameta.allure.Step;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.nio.charset.StandardCharsets;
import java.time.Duration;

public final class UiClient {

    private static final HttpClient HTTP = HttpClient.newBuilder()
            .connectTimeout(Duration.ofSeconds(5))
            .build();

    private UiClient() {
    }

    @Step("GET UI {uiUrl}status")
    public static HubClient.HubStatusResponse fetchStatus(String uiUrl) {
        return HubClient.fetchStatus(uiUrl);
    }

    @Step("Read first SSE event from {uiUrl}events")
    public static JsonObject readFirstSseEvent(String uiUrl) {
        var url = uiUrl.endsWith("/") ? uiUrl + "events" : uiUrl + "/events";
        var request = HttpRequest.newBuilder(URI.create(url))
                .timeout(Duration.ofSeconds(20))
                .header("Accept", "text/event-stream")
                .GET()
                .build();
        try {
            var response = HTTP.send(request, HttpResponse.BodyHandlers.ofInputStream());
            if (response.statusCode() != 200) {
                throw new IllegalStateException("SSE GET " + url + " returned HTTP " + response.statusCode());
            }
            try (var reader = new BufferedReader(new InputStreamReader(response.body(), StandardCharsets.UTF_8))) {
                var deadline = System.currentTimeMillis() + 15_000;
                String line;
                while ((line = reader.readLine()) != null) {
                    if (System.currentTimeMillis() > deadline) {
                        throw new IllegalStateException("SSE timeout waiting for first event");
                    }
                    if (line.startsWith("data:")) {
                        var payload = line.substring("data:".length()).trim();
                        JsonElement parsed = JsonParser.parseString(payload);
                        if (!parsed.isJsonObject()) {
                            throw new IllegalStateException("Expected SSE JSON object, got: " + payload);
                        }
                        return parsed.getAsJsonObject();
                    }
                }
            }
            throw new IllegalStateException("SSE stream ended without data event");
        } catch (IOException | InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new IllegalStateException("Failed to read SSE from " + url + ": " + e.getMessage(), e);
        }
    }
}
