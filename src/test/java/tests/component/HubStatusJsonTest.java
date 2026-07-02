package tests.component;

import annotations.Layer;
import api.hub.HubStatus;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;

@Layer("component")
@DisplayName("Hub status fixture")
class HubStatusJsonTest {

    @Test
    @DisplayName("parses idle hub status counters")
    void parsesIdleHubStatus() {
        var status = FixtureJson.parse("fixtures/hub/status-idle.json", HubStatus.class);
        assertEquals(5, status.total());
        assertEquals(0, status.used());
        assertNotNull(status.browsers());
    }
}
