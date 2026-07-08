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

import static io.qameta.allure.Allure.step;

@Layer("api")
@Component("selenoid")
@Epic("selenoid")
@Feature("Hub logs WebSocket")
@DisplayName("Hub session logs WebSocket API")
@ResourceLock(value = "hubSessions", mode = ResourceAccessMode.READ_WRITE)
class HubLogsSessionApiTests extends ApiTestBase {

    @Test
    @Tag("api")
    @Tag("positive")
    @DisplayName("WebSocket /logs/{sessionId} opens for active session")
    void sessionLogsWebSocketOpens() {
        var sessionId = step("Create hub session", () -> HubSessionApi.create());
        try {
            step("Open logs WebSocket", () -> {
                try {
                    HubWebSocketApi.openBriefly("/logs/" + sessionId, Duration.ofSeconds(15));
                } catch (Exception e) {
                    throw new RuntimeException(e);
                }
            });
        } finally {
            step("Delete hub session", () -> HubSessionApi.delete(sessionId));
        }
    }
}
