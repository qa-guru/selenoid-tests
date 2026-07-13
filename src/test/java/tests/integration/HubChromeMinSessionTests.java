package tests.integration;

import annotations.Component;
import annotations.Layer;
import api.hub.HubSessionApi;
import api.hub.HubStatusApi;
import config.ConfigReader;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import io.qameta.allure.Story;
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
@Story("WebDriver session (min)")
@DisplayName("WebDriver hub session (chrome-min)")
@ResourceLock(value = "hubSessions", mode = ResourceAccessMode.READ_WRITE)
class HubChromeMinSessionTests {

    @Test
    @Tag("integration")
    @Tag("min")
    @Tag("positive")
    @DisplayName("Remote Chrome min WebDriver session starts browser node and completes")
    void remoteChromeMinSessionStartsAndCompletes() {
        var usedBefore = step("Snapshot hub /status used counter", () ->
                HubStatusApi.fetch().used());

        var sessionId = step("Create Chrome min hub session", () ->
                HubSessionApi.create("chrome", ConfigReader.testConfig.chromeMinVersion()));

        try {
            step("Verify hub reports active session", () ->
                    assertEquals(usedBefore + 1, HubStatusApi.fetch().used(),
                            "Hub /status.used should increment while Chrome min session is open"));

            step("Verify session id is assigned", () ->
                    assertFalse(sessionId.isBlank(), "Chrome min session id should not be blank"));
        } finally {
            step("Delete Chrome min session", () -> HubSessionApi.delete(sessionId));
        }

        step("Verify hub released session", () ->
                assertEquals(usedBefore, HubStatusApi.fetch().used(),
                        "Hub /status.used should return to baseline after Chrome min session ends"));
    }
}
