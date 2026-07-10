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
@Feature("WebDriver session (min)")
@DisplayName("WebDriver hub session (msedge-min)")
@ResourceLock(value = "hubSessions", mode = ResourceAccessMode.READ_WRITE)
class HubMsedgeMinSessionTests {

    private static final String MSEDGE_MIN_VERSION = "145.0-min";

    @Test
    @Tag("integration")
    @Tag("min")
    @Tag("positive")
    @DisplayName("Remote Edge min WebDriver session starts browser node and completes")
    void remoteMsedgeMinSessionStartsAndCompletes() {
        var usedBefore = step("Snapshot hub /status used counter", () ->
                HubStatusApi.fetch().used());

        var sessionId = step("Create Edge min hub session", () ->
                HubSessionApi.create("msedge", MSEDGE_MIN_VERSION));

        try {
            step("Verify hub reports active session", () ->
                    assertEquals(usedBefore + 1, HubStatusApi.fetch().used(),
                            "Hub /status.used should increment while Edge min session is open"));

            step("Verify session id is assigned", () ->
                    assertFalse(sessionId.isBlank(), "Edge min session id should not be blank"));
        } finally {
            step("Delete Edge min session", () -> HubSessionApi.delete(sessionId));
        }

        step("Verify hub released session", () ->
                assertEquals(usedBefore, HubStatusApi.fetch().used(),
                        "Hub /status.used should return to baseline after Edge min session ends"));
    }
}
