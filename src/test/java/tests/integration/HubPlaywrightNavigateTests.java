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

import static io.qameta.allure.Allure.step;
import static org.junit.jupiter.api.Assertions.assertEquals;

@Layer("integration")
@Component("playwright-image")
@Epic("playwright-image")
@Feature("Playwright WS session")
@DisplayName("Playwright navigate integration")
class HubPlaywrightNavigateTests {

    @Test
    @Tag("integration")
    @Tag("positive")
    @DisplayName("Playwright session navigates to example.com")
    void playwrightSessionNavigatesToExampleDomain() {
        var usedBefore = step("Snapshot hub used counter", () -> HubStatusApi.fetch().used());
        try (var playwright = Playwright.create()) {
            var browser = step("Connect Playwright via hub WS endpoint", () ->
                    PlaywrightSessionApi.connect(playwright));
            try (var page = browser.newPage()) {
                step("Navigate to example.com", () -> page.navigate("https://example.com/"));
                step("Verify page title", () ->
                        assertEquals("Example Domain", page.title()));
            } finally {
                step("Close remote session", () -> PlaywrightSessionApi.close(browser));
            }
        }
        step("Verify hub released session", () ->
                assertEquals(usedBefore, HubStatusApi.fetch().used()));
    }
}
