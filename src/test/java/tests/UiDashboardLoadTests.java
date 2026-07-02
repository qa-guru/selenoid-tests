package tests;

import annotations.Layer;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import static com.codeborne.selenide.Selenide.open;
import static com.codeborne.selenide.WebDriverRunner.url;
import static io.qameta.allure.Allure.step;
import static org.junit.jupiter.api.Assertions.assertTrue;

@Layer("e2e")
@Epic("selenoid-ui")
@Feature("UI dashboard")
@DisplayName("UI dashboard load")
class UiDashboardLoadTests extends UiTestBase {

    @Test
    @Tag("smoke")
    @Tag("positive")
    @DisplayName("Dashboard opens root URL")
    void dashboardOpensRootUrl() {
        step("Open dashboard", () -> uiDashboard.openPage());
        step("Verify root URL is loaded", () -> assertTrue(url().endsWith("/") || url().endsWith(":8080/")));
    }
}
