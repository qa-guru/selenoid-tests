package tests.api;

import annotations.Component;
import annotations.Layer;
import api.UiApiTestBase;
import api.hub.HubSessionApi;
import api.ui.UiWebSocketApi;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.parallel.ResourceAccessMode;
import org.junit.jupiter.api.parallel.ResourceLock;

import java.time.Duration;

import static io.qameta.allure.Allure.step;

@Layer("api")
@Component("selenoid-ui")
@Epic("selenoid-ui")
@Feature("UI logs WebSocket proxy")
@DisplayName("UI /ws/logs WebSocket API")
@ResourceLock(value = "hubSessions", mode = ResourceAccessMode.READ_WRITE)
class UiLogsWsApiTests extends UiApiTestBase {

    @Test
    @Tag("api")
    @Tag("positive")
    @DisplayName("WebSocket /ws/logs/{sessionId} opens via UI proxy")
    void uiLogsWebSocketOpens() {
        var sessionId = step("Create hub session", () -> HubSessionApi.create());
        try {
            step("Open UI logs WebSocket", () -> {
                try {
                    UiWebSocketApi.openBriefly("/ws/logs/" + sessionId, Duration.ofSeconds(15));
                } catch (Exception e) {
                    throw new RuntimeException(e);
                }
            });
        } finally {
            step("Delete hub session", () -> HubSessionApi.delete(sessionId));
        }
    }
}
