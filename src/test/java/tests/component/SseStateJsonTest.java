package tests.component;

import annotations.Component;
import annotations.Layer;
import api.ui.SseHubEvent;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertTrue;

@Component("selenoid-ui")
@Layer("component")
@DisplayName("SSE state fixture")
class SseStateJsonTest {

    @Test
    @DisplayName("parses SSE payload with hub state")
    void parsesSseStatePayload() {
        var event = FixtureJson.parse("fixtures/sse/state.json", SseHubEvent.class);
        assertTrue(event.hasState());
    }
}
