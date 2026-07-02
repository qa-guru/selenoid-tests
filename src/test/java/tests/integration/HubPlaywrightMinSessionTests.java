package tests.integration;

import annotations.Component;
import annotations.Layer;
import api.hub.HubStatusApi;
import api.hub.PlaywrightSessionApi;
import com.microsoft.playwright.Playwright;
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
import static org.junit.jupiter.api.Assertions.assertTrue;

@Layer("integration")
@Component("playwright-image")
@Epic("playwright-image")
@Feature("Playwright WS session (min)")
@DisplayName("Playwright hub WS session (min)")
@ResourceLock(value = "hubSessions", mode = ResourceAccessMode.READ_WRITE)
class HubPlaywrightMinSessionTests {

    @Test
    @Tag("min")
    @Tag("positive")
    @DisplayName("Remote Playwright min WS session starts browser node and completes")
    void remotePlaywrightMinSessionStartsAndCompletes() {
        var usedBefore = step("Snapshot hub /status used counter", () ->
                HubStatusApi.fetch().used());

        try (var playwright = Playwright.create()) {
            var browser = step("Connect Playwright via hub WS endpoint (min)", () ->
                    PlaywrightSessionApi.connect(playwright));

            step("Verify remote browser is connected", () ->
                    assertTrue(browser.isConnected(), "Playwright min browser node should be connected"));

            step("Verify hub reports active session", () ->
                    assertEquals(usedBefore + 1, HubStatusApi.fetch().used(),
                            "Hub /status.used should increment while Playwright min session is open"));

            step("Close remote session", () -> PlaywrightSessionApi.close(browser));

            step("Verify browser is disconnected after close", () ->
                    assertFalse(browser.isConnected(), "Browser should disconnect after close"));
        }

        step("Verify hub released session", () ->
                assertEquals(usedBefore, HubStatusApi.fetch().used(),
                        "Hub /status.used should return to baseline after min session ends"));
    }
}
