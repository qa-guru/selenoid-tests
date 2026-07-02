package tests.component;

import annotations.Layer;
import api.ui.SseHubEvent;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertTrue;

@Layer("component")
@DisplayName("SSE errors fixture")
class SseErrorsJsonTest {

    @Test
    @DisplayName("parses SSE payload with errors list")
    void parsesSseErrorsPayload() {
        var event = FixtureJson.parse("fixtures/sse/errors.json", SseHubEvent.class);
        assertTrue(event.hasErrors());
    }
}
