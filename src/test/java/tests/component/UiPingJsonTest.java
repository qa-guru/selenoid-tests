package tests.component;

import annotations.Layer;
import api.ui.UiPingResponse;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;

@Layer("component")
@DisplayName("UI ping fixture")
class UiPingJsonTest {

    @Test
    @DisplayName("parses uptime and version fields")
    void parsesPingFields() {
        var ping = FixtureJson.parse("fixtures/ui/ping.json", UiPingResponse.class);
        assertEquals("1h2m3s", ping.uptime());
        assertEquals("1.0.0-test", ping.version());
    }
}
