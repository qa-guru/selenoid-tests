package tests;

import annotations.Component;
import annotations.Layer;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import io.qameta.allure.Story;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import static com.microsoft.playwright.assertions.PlaywrightAssertions.assertThat;
import static io.qameta.allure.Allure.step;
import static org.junit.jupiter.api.Assertions.assertFalse;

@Layer("e2e")
@Component("playwright-image")
@Epic("playwright-image")
@Feature("Playwright WS session (min)")
@Story("Playwright WS session (min)")
@DisplayName("Playwright WS session (min)")
class HubPlaywrightMinSessionTests extends PlaywrightTestBase {

    @Test
    @Tag("min")
    @Tag("positive")
    @DisplayName("Remote Playwright min WS session opens example.com")
    void remotePlaywrightMinSessionOpensExampleDomain() {
        step("Navigate to smoke URL", () -> page.navigate(config.smokeUrl()));

        step("Verify page title", () ->
                assertThat(page).hasTitle("Example Domain"));

        step("Verify heading text", () ->
                assertThat(page.locator("h1")).hasText("Example Domain"));

        step("Verify browser is connected", () ->
                assertFalse(page.isClosed()));
    }
}
