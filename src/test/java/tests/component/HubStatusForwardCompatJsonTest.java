package tests.component;

import annotations.Component;
import annotations.Layer;
import api.hub.HubStatus;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;

@Component("selenoid")
@Layer("component")
@DisplayName("Hub status forward compatibility")
class HubStatusForwardCompatJsonTest {

    @Test
    @DisplayName("ignores unknown JSON fields")
    void ignoresUnknownFields() {
        var status = FixtureJson.parse("fixtures/hub/status-unknown-field.json", HubStatus.class);
        assertEquals(3, status.total());
    }
}
