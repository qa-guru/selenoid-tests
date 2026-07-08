package tests.component;

import annotations.Component;
import annotations.Layer;
import api.hub.HubStatusParser;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;

@Component("selenoid")
@Layer("component")
@DisplayName("Hub status parser")
class HubStatusParserTest {

    @Test
    @DisplayName("parses flat hub /status JSON")
    void parsesFlatHubStatus() {
        var status = HubStatusParser.parseJson(FixtureJson.load("fixtures/hub/status-idle.json"));
        assertEquals(0, status.used());
        assertEquals(5, status.total());
    }

    @Test
    @DisplayName("parses UI-shaped /status JSON with .state wrapper")
    void parsesUiWrappedStatus() {
        var status = HubStatusParser.parseJson(FixtureJson.load("fixtures/ui/status.json"));
        assertEquals(1, status.used());
        assertEquals(5, status.total());
    }
}
