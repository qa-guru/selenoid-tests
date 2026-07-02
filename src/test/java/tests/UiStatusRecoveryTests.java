package tests;

import annotations.Component;
import annotations.Layer;
import allure.Attachments;
import helpers.StackHelper;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.Timeout;
import org.junit.jupiter.api.parallel.Execution;
import org.junit.jupiter.api.parallel.ExecutionMode;

import static com.codeborne.selenide.Selenide.closeWebDriver;
import static io.qameta.allure.Allure.step;

@Layer("integration")
@Component("selenoid-ui")
@Epic("selenoid-ui")
@Feature("UI status recovery")
@DisplayName("UI status recovery")
@Tag("resilience")
@Tag("integration")
@Tag("local-only")
@Execution(ExecutionMode.SAME_THREAD)
class UiStatusRecoveryTests extends UiTestBase {

    @Test
    @Tag("positive")
    @Timeout(120)
    @DisplayName("SELENOID recovers after hub restart")
    void selenoidRecoversAfterHubRestart() throws Exception {
        step("Baseline CONNECTED", () -> uiDashboard.openPage().shouldBeConnected());
        attachScreenshot("Before hub kill");

        step("Stop hub", () -> {
            StackHelper.killHub();
            StackHelper.waitForHubDown(15_000);
        });
        uiDashboard.openPage().shouldBeDegraded();
        attachScreenshot("Hub down");

        step("Restart hub", () -> {
            StackHelper.startHubDetached();
            StackHelper.waitForHubReady(30_000);
        });

        step("Fresh browser after hub recovery", () -> closeWebDriver());

        step("Both indicators recover within UI retry window", () -> {
            uiDashboard.openPage();
            uiDashboard.shouldSseBeConnected();
            uiDashboard.shouldSelenoidBeConnected();
            uiDashboard.shouldStayConnected(5_000);
        });

        attachScreenshot("Recovered");
    }

    private void attachScreenshot(String name) {
        Attachments.screenshot(name);
    }
}
