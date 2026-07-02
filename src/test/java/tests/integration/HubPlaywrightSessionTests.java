package tests.integration;

import annotations.Layer;
import api.hub.HubStatusApi;
import api.hub.PlaywrightSessionApi;
import com.microsoft.playwright.Playwright;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import static io.qameta.allure.Allure.step;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertTrue;

@Layer("integration")
@Epic("playwright-image")
@Feature("Playwright WS session")
@DisplayName("Playwright hub WS session")
class HubPlaywrightSessionTests {

    @Test
    @Tag("integration")
    @Tag("positive")
    @DisplayName("Remote Playwright WS session starts browser node and completes")
    void remotePlaywrightSessionStartsAndCompletes() {
        var usedBefore = step("Snapshot hub /status used counter", () ->
                HubStatusApi.fetch().used());

        try (var playwright = Playwright.create()) {
            var browser = step("Connect Playwright via hub WS endpoint", () ->
                    PlaywrightSessionApi.connect(playwright));

            step("Verify remote browser is connected", () ->
                    assertTrue(browser.isConnected(), "Playwright browser node should be connected"));

            step("Verify hub reports active session", () ->
                    assertEquals(usedBefore + 1, HubStatusApi.fetch().used(),
                            "Hub /status.used should increment while Playwright session is open"));

            step("Close remote session", () -> PlaywrightSessionApi.close(browser));

            step("Verify browser is disconnected after close", () ->
                    assertFalse(browser.isConnected(), "Browser should disconnect after close"));
        }

        step("Verify hub released session", () ->
                assertEquals(usedBefore, HubStatusApi.fetch().used(),
                        "Hub /status.used should return to baseline after session ends"));
    }
}
