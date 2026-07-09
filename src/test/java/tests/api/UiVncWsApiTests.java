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
import java.util.Map;

import static io.qameta.allure.Allure.step;

@Layer("api")
@Component("selenoid-ui")
@Epic("selenoid-ui")
@Feature("UI VNC WebSocket proxy")
@DisplayName("UI /ws/vnc WebSocket API")
@ResourceLock(value = "hubSessions", mode = ResourceAccessMode.READ_WRITE)
class UiVncWsApiTests extends UiApiTestBase {

    @Test
    @Tag("api")
    @Tag("positive")
    @DisplayName("WebSocket /ws/vnc/{sessionId} opens via UI when enableVNC=true")
    void uiVncWebSocketOpens() {
        var sessionId = step("Create hub session with VNC", () ->
                HubSessionApi.createWithSelenoidOptions(config, Map.of("enableVNC", true)));
        try {
            step("Open UI VNC WebSocket", () -> {
                try {
                    UiWebSocketApi.openBriefly("/ws/vnc/" + sessionId, Duration.ofSeconds(15));
                } catch (Exception e) {
                    throw new RuntimeException(e);
                }
            });
        } finally {
            step("Delete hub session", () -> HubSessionApi.delete(sessionId));
        }
    }
}
