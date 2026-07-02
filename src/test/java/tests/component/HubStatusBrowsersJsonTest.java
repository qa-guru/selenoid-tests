package tests.component;

import annotations.Layer;
import api.hub.HubStatus;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertTrue;

@Layer("component")
@DisplayName("Hub status browsers fixture")
class HubStatusBrowsersJsonTest {

    @Test
    @DisplayName("parses browsers map with version entries")
    void parsesBrowsersMap() {
        var status = FixtureJson.parse("fixtures/hub/status-with-browsers.json", HubStatus.class);
        assertTrue(status.browsers().containsKey("chrome"));
        assertFalse(status.browsers().isEmpty());
    }
}
