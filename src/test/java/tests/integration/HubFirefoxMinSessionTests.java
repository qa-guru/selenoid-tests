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
@DisplayName("WebDriver hub session (firefox-min)")
@ResourceLock(value = "hubSessions", mode = ResourceAccessMode.READ_WRITE)
class HubFirefoxMinSessionTests {

    @Test
    @Tag("integration")
    @Tag("min")
    @Tag("positive")
    @DisplayName("Remote Firefox min WebDriver session starts browser node and completes")
    void remoteFirefoxMinSessionStartsAndCompletes() {
        var usedBefore = step("Snapshot hub /status used counter", () ->
                HubStatusApi.fetch().used());

        var sessionId = step("Create Firefox min hub session", () ->
                HubSessionApi.create("firefox", ConfigReader.testConfig.firefoxMinVersion()));

        try {
            step("Verify hub reports active session", () ->
                    assertEquals(usedBefore + 1, HubStatusApi.fetch().used(),
                            "Hub /status.used should increment while Firefox min session is open"));

            step("Verify session id is assigned", () ->
                    assertFalse(sessionId.isBlank(), "Firefox min session id should not be blank"));
        } finally {
            step("Delete Firefox min session", () -> HubSessionApi.delete(sessionId));
        }

        step("Verify hub released session", () ->
                assertEquals(usedBefore, HubStatusApi.fetch().used(),
                        "Hub /status.used should return to baseline after Firefox min session ends"));
    }
}
