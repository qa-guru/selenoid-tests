package tests.api;

import annotations.Component;
import annotations.Layer;
import api.ApiTestBase;
import api.hub.HubSessionApi;
import api.hub.HubWebSocketApi;
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
import static org.junit.jupiter.api.Assertions.assertTrue;

@Layer("api")
@Component("selenoid")
@Epic("selenoid")
@Feature("Hub VNC WebSocket")
@DisplayName("Hub session VNC WebSocket API")
@ResourceLock(value = "hubSessions", mode = ResourceAccessMode.READ_WRITE)
class HubVncSessionApiTests extends ApiTestBase {

    @Test
    @Tag("api")
    @Tag("positive")
    @DisplayName("WebSocket /vnc/{sessionId} opens when enableVNC=true")
    void sessionVncWebSocketOpens() {
        var sessionId = step("Create hub session with VNC", () ->
                HubSessionApi.createWithSelenoidOptions(config, Map.of("enableVNC", true)));
        try {
            step("Open VNC WebSocket", () -> {
                try {
                    HubWebSocketApi.openBriefly("/vnc/" + sessionId, Duration.ofSeconds(15));
                } catch (Exception e) {
                    throw new RuntimeException(e);
                }
            });
        } finally {
            step("Delete hub session", () -> HubSessionApi.delete(sessionId));
        }
    }
}
