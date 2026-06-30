package helpers;

import com.google.gson.JsonElement;
import com.google.gson.JsonObject;
import com.google.gson.JsonParser;
import io.qameta.allure.Step;

import java.io.IOException;
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.time.Duration;

public final class HubClient {

    private static final HttpClient HTTP = HttpClient.newBuilder()
            .connectTimeout(Duration.ofSeconds(5))
            .build();

    private HubClient() {
    }

    @Step("GET {hubUrl}status")
    public static HubStatusResponse fetchStatus(String hubUrl) {
        var url = hubUrl.endsWith("/") ? hubUrl + "status" : hubUrl + "/status";
        var request = HttpRequest.newBuilder(URI.create(url))
                .timeout(Duration.ofSeconds(10))
                .GET()
                .build();
        try {
            var response = HTTP.send(request, HttpResponse.BodyHandlers.ofString());
            return new HubStatusResponse(response.statusCode(), response.body());
        } catch (IOException | InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new IllegalStateException("Failed to GET " + url + ": " + e.getMessage(), e);
        }
    }

    public record HubStatusResponse(int httpStatus, String body) {

        public JsonObject json() {
            JsonElement parsed = JsonParser.parseString(body);
            if (!parsed.isJsonObject()) {
                throw new IllegalStateException("Expected JSON object, got: " + body);
            }
            return parsed.getAsJsonObject();
        }
    }
}
