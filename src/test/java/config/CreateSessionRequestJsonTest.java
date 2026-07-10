package config;

import annotations.Component;
import annotations.Layer;
import api.hub.HubSessionApi;
import com.google.gson.Gson;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertTrue;

@Component("selenoid")
@Layer("unit")
@DisplayName("CreateSessionRequest JSON")
class CreateSessionRequestJsonTest {

    private static final Gson GSON = new Gson();

    @Test
    @DisplayName("session body includes docker-safe chrome args")
    void sessionBodyIncludesDockerSafeChromeArgs() {
        var json = GSON.toJson(HubSessionApi.createSessionBody("chrome", "148.0-min"));
        assertTrue(json.contains("browserName"));
        assertTrue(json.contains("chrome"));
        assertTrue(json.contains("browserVersion"));
        assertTrue(json.contains("148.0-min"));
        assertTrue(json.contains("goog:chromeOptions"));
        assertTrue(json.contains("no-sandbox"));
        assertTrue(json.contains("disable-dev-shm-usage"));
    }
}
