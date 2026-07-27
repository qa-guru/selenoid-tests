package tests;

import annotations.Component;
import annotations.Layer;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import io.qameta.allure.Story;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import static com.codeborne.selenide.WebDriverRunner.url;
import static io.qameta.allure.Allure.step;
import static org.junit.jupiter.api.Assertions.assertTrue;

@Layer("e2e")
@Component("selenoid-ui")
@Epic("selenoid-ui")
@Feature("UI dashboard")
@Story("UI dashboard load")
@DisplayName("UI dashboard load")
class UiDashboardLoadTests extends UiTestBase {

    @Test
    @Tag("smoke")
    @Tag("positive")
    @DisplayName("Dashboard opens root URL")
    void dashboardOpensRootUrl() {
        step("Open dashboard", () -> uiDashboard.openPage());
        // v3 HashRouter: `/` redirects to `#/statistics` (legacy root `/` or `:8080/` still OK).
        step("Verify Statistics route (or root) is loaded", () -> {
            var current = url();
            assertTrue(
                    current.contains("/#/statistics")
                            || current.endsWith("/")
                            || current.endsWith(":8080/"),
                    () -> "Expected Statistics hash route or root, got: " + current);
        });
    }
}
