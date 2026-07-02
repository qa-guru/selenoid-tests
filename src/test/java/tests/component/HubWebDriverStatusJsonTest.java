package tests.component;

import annotations.Component;
import annotations.Layer;
import api.hub.HubWebDriverStatus;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertTrue;

@Component("selenoid")
@Layer("component")
@DisplayName("Hub WebDriver status fixture")
class HubWebDriverStatusJsonTest {

    @Test
    @DisplayName("parses ready WebDriver status payload")
    void parsesReadyWebDriverStatus() {
        var status = FixtureJson.parse("fixtures/hub/wd-status-ready.json", HubWebDriverStatus.class);
        assertTrue(status.value().ready());
        assertEquals("Selenoid ready", status.value().message());
    }
}
