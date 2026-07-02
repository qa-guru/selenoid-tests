package tests;

import annotations.Component;
import annotations.Layer;
import io.qameta.allure.AllureId;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import static com.microsoft.playwright.assertions.PlaywrightAssertions.assertThat;
import static io.qameta.allure.Allure.step;
import static org.junit.jupiter.api.Assertions.assertFalse;

@Layer("e2e")
@Component("playwright-image")
@Epic("playwright-image")
@Feature("Playwright WS session")
@DisplayName("Playwright WS session")
class HubPlaywrightSessionTests extends PlaywrightTestBase {

    @Test
    @AllureId("45342")
    @Tag("playwright")
    @Tag("smoke")
    @Tag("positive")
    @DisplayName("Remote Playwright WS session opens example.com")
    void remotePlaywrightSessionOpensExampleDomain() {
        step("Navigate to smoke URL", () -> page.navigate(config.smokeUrl()));

        step("Verify page title", () ->
                assertThat(page).hasTitle("Example Domain"));

        step("Verify heading text", () ->
                assertThat(page.locator("h1")).hasText("Example Domain"));

        step("Verify browser is connected", () ->
                assertFalse(page.isClosed()));
    }
}
