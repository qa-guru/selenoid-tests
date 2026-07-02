package tests.component;

import annotations.Layer;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;

@Layer("component")
@DisplayName("Session create fixture")
class SessionCreateJsonTest {

    @Test
    @DisplayName("extracts sessionId from create session response")
    void extractsSessionId() {
        var sessionId = FixtureJson.load("fixtures/hub/session-create.json");
        assertEquals("abc123-session",
                io.restassured.path.json.JsonPath.from(sessionId).getString("value.sessionId"));
    }
}
