package tests.integration;

import annotations.Component;
import annotations.Layer;
import api.hub.HubSessionApi;
import api.hub.HubStatusApi;
import config.ConfigReader;
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
@DisplayName("WebDriver hub session (chrome)")
@ResourceLock(value = "hubSessions", mode = ResourceAccessMode.READ_WRITE)
class HubChromeSessionIntegrationTests {

    @Test
    @Tag("integration")
    @Tag("positive")
    @DisplayName("Remote Chrome WebDriver session starts browser node and completes")
    void remoteChromeSessionStartsAndCompletes() {
        var browserVersion = ConfigReader.testConfig.browserVersion();
        var usedBefore = step("Snapshot hub /status used counter", () ->
                HubStatusApi.fetch().used());

        var sessionId = step("Create Chrome hub session", () ->
                HubSessionApi.create("chrome", browserVersion));

        try {
            step("Verify hub reports active session", () ->
                    assertEquals(usedBefore + 1, HubStatusApi.fetch().used(),
                            "Hub /status.used should increment while Chrome session is open"));

            step("Verify session id is assigned", () ->
                    assertFalse(sessionId.isBlank(), "Chrome session id should not be blank"));
        } finally {
            step("Delete Chrome session", () -> HubSessionApi.delete(sessionId));
        }

        step("Verify hub released session", () ->
                assertEquals(usedBefore, HubStatusApi.fetch().used(),
                        "Hub /status.used should return to baseline after Chrome session ends"));
    }
}
