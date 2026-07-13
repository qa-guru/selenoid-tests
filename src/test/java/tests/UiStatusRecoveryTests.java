package tests;

import annotations.Component;
import annotations.Layer;
import allure.Attachments;
import api.ui.UiStatusApi;
import helpers.StackHelper;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import io.qameta.allure.Story;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.Timeout;
import org.junit.jupiter.api.parallel.Execution;
import org.junit.jupiter.api.parallel.ExecutionMode;

import static com.codeborne.selenide.Selenide.closeWebDriver;
import static io.qameta.allure.Allure.step;
import static org.junit.jupiter.api.Assertions.assertFalse;

@Layer("integration")
@Component("selenoid-ui")
@Epic("selenoid-ui")
@Feature("UI status recovery")
@Story("UI status recovery")
@DisplayName("UI status recovery")
@Tag("resilience")
@Execution(ExecutionMode.SAME_THREAD)
class UiStatusRecoveryTests extends UiTestBase {

    @Test
    @Tag("positive")
    @Timeout(120)
    @DisplayName("SELENOID recovers after hub restart")
    void selenoidRecoversAfterHubRestart() throws Exception {
        step("Ensure hub and UI are up", () -> StackHelper.ensureStackRunning());
        step("Baseline CONNECTED", () -> uiDashboard.openPage().shouldBeConnected());
        attachScreenshot("Before hub kill");

        step("Stop hub", () -> {
            StackHelper.killHub();
            StackHelper.waitForHubDown(15_000);
        });
        step("Close browser after hub stop", () -> closeWebDriver());
        step("UI /status reports hub error", () -> {
            var response = UiStatusApi.fetchWhenHubUnavailable();
            assertFalse(response.errors().isEmpty());
        });

        step("Restart hub", () -> {
            StackHelper.startHubDetached();
            StackHelper.waitForHubReady(30_000);
        });

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
