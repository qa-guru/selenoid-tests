package tests.integration;

import annotations.Component;
import annotations.Layer;
import api.hub.HubStatusApi;
import api.hub.PlaywrightSessionApi;
import com.microsoft.playwright.Playwright;
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
import static org.junit.jupiter.api.Assertions.assertTrue;

@Layer("integration")
@Component("playwright-image")
@Epic("playwright-image")
@Feature("Playwright Firefox WS session")
@Story("Playwright Firefox WS session")
@DisplayName("Playwright Firefox hub WS session")
@ResourceLock(value = "hubSessions", mode = ResourceAccessMode.READ_WRITE)
class HubPlaywrightFirefoxSessionTests {

    @Test
    @Tag("integration")
    @Tag("playwright")
    @DisplayName("Remote Playwright Firefox WS session connects via hub")
    void remoteFirefoxSessionConnects() {
        var endpoint = ConfigReader.resolvePlaywrightWsEndpoint(ConfigReader.testConfig, "playwright-firefox");
        var usedBefore = step("Snapshot hub /status used counter", () ->
                HubStatusApi.fetch().used());

        try (var playwright = Playwright.create()) {
            var browser = step("Connect Playwright Firefox via hub WS", () ->
                    PlaywrightSessionApi.connect(playwright, endpoint));

            step("Verify remote browser is connected", () ->
                    assertTrue(browser.isConnected(), "Firefox browser node should be connected"));

            step("Verify hub reports active session", () ->
                    assertEquals(usedBefore + 1, HubStatusApi.fetch().used()));

            step("Close remote session", () -> PlaywrightSessionApi.close(browser));

            step("Verify browser is disconnected", () ->
                    assertFalse(browser.isConnected()));
        }

        step("Verify hub released session", () ->
                assertEquals(usedBefore, HubStatusApi.fetch().used()));
    }
}
