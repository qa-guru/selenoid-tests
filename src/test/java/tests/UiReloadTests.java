package tests;

import annotations.Layer;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import static io.qameta.allure.Allure.step;

@Layer("e2e")
@Epic("selenoid-ui")
@Feature("UI dashboard")
@DisplayName("UI dashboard reload")
class UiReloadTests extends UiTestBase {

    @Test
    @Tag("smoke")
    @Tag("positive")
    @DisplayName("Dashboard reload keeps CONNECTED indicators")
    void dashboardReloadKeepsConnected() {
        step("Open dashboard and wait for CONNECTED", () ->
                uiDashboard.openPage().shouldBeConnected());
        step("Reload dashboard", () -> uiDashboard.reloadPage().shouldBeConnected());
    }
}
