package tests.integration;

import annotations.Component;
import annotations.Layer;
import api.hub.HubSessionApi;
import api.ui.SseStreamApi;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import io.qameta.allure.Story;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import org.junit.jupiter.api.parallel.ResourceAccessMode;
import org.junit.jupiter.api.parallel.ResourceLock;

import static io.qameta.allure.Allure.step;
import static org.junit.jupiter.api.Assertions.assertTrue;

@Layer("integration")
@Component("selenoid-ui")
@Epic("selenoid-ui")
@Feature("UI SSE")
@Story("UI SSE")
@DisplayName("UI SSE with hub session")
@ResourceLock(value = "hubSessions", mode = ResourceAccessMode.READ_WRITE)
class UiSseWithSessionTests {

    @Test
    @Tag("integration")
    @Tag("positive")
    @DisplayName("SSE payload remains parseable while hub session is active")
    void ssePayloadParseableWithActiveSession() {
        var sessionId = step("Create hub session", () -> HubSessionApi.create());
        try {
            var event = step("Read SSE event", () -> SseStreamApi.readFirstEvent());
            step("Verify SSE payload has state or errors", () ->
                    assertTrue(event.hasState() || event.hasErrors()));
        } finally {
            step("Delete hub session", () -> HubSessionApi.delete(sessionId));
        }
    }
}
