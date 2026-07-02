package config;

import annotations.Layer;
import api.hub.CreateSessionCapabilities;
import api.hub.CreateSessionRequest;
import api.hub.SessionAlwaysMatch;
import com.google.gson.Gson;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertTrue;

@Layer("unit")
@DisplayName("CreateSessionRequest JSON")
class CreateSessionRequestJsonTest {

    private static final Gson GSON = new Gson();

    @Test
    @DisplayName("serializes alwaysMatch browser capabilities")
    void serializesAlwaysMatchCapabilities() {
        var request = new CreateSessionRequest(
                new CreateSessionCapabilities(new SessionAlwaysMatch("chrome", "148.0")));
        var json = GSON.toJson(request);
        assertTrue(json.contains("alwaysMatch"));
        assertTrue(json.contains("browserName"));
        assertTrue(json.contains("chrome"));
        assertTrue(json.contains("browserVersion"));
        assertTrue(json.contains("148.0"));
    }
}
