package tests.integration;

import annotations.Component;
import annotations.Layer;
import api.hub.HubSessionApi;
import api.hub.HubStatusApi;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.parallel.ResourceAccessMode;
import org.junit.jupiter.api.parallel.ResourceLock;

import static io.qameta.allure.Allure.step;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;

@Layer("integration")
@Component("webdriver-image")
@Epic("webdriver-image")
@Feature("WebDriver session")
@DisplayName("WebDriver hub session (msedge)")
@ResourceLock(value = "hubSessions", mode = ResourceAccessMode.READ_WRITE)
class HubMsedgeSessionIntegrationTests {

    private static final String MSEDGE_WARM_VERSION = "145.0";

    @Test
    @Tag("integration")
    @Tag("positive")
    @DisplayName("Remote Edge WebDriver session starts browser node and completes")
    void remoteMsedgeSessionStartsAndCompletes() {
        var usedBefore = step("Snapshot hub /status used counter", () ->
                HubStatusApi.fetch().used());

        var sessionId = step("Create Edge hub session", () ->
                HubSessionApi.create("msedge", MSEDGE_WARM_VERSION));

        try {
            step("Verify hub reports active session", () ->
                    assertEquals(usedBefore + 1, HubStatusApi.fetch().used(),
                            "Hub /status.used should increment while Edge session is open"));

            step("Verify session id is assigned", () ->
                    assertFalse(sessionId.isBlank(), "Edge session id should not be blank"));
        } finally {
            step("Delete Edge session", () -> HubSessionApi.delete(sessionId));
        }

        step("Verify hub released session", () ->
                assertEquals(usedBefore, HubStatusApi.fetch().used(),
                        "Hub /status.used should return to baseline after Edge session ends"));
    }
}
