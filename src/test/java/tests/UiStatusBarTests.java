package tests;

import annotations.Component;
import annotations.Layer;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import io.qameta.allure.Story;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import static com.codeborne.selenide.Selenide.$;
import static io.qameta.allure.Allure.step;

@Layer("e2e")
@Component("selenoid-ui")
@Epic("selenoid-ui")
@Feature("UI status bar")
@Story("UI status bar")
@DisplayName("UI status bar")
class UiStatusBarTests extends UiTestBase {

    @Test
    @Tag("smoke")
    @Tag("positive")
    @DisplayName("SSE and SELENOID stay CONNECTED on stable stack")
    void statusBarStaysConnected() throws InterruptedException {
        step("Open dashboard and wait for CONNECTED", () ->
                uiDashboard.openPage().shouldBeConnected());

        step("Verify status indicators show connected StatusTile styling", () -> {
            $("#sse-status").shouldHave(com.codeborne.selenide.Condition.cssClass("status-tile--connected"));
            $("#selenoid-status").shouldHave(com.codeborne.selenide.Condition.cssClass("status-tile--connected"));
        });

        step("Keep CONNECTED stable for 5 seconds", () ->
                uiDashboard.shouldStayConnected(5_000));
    }
}
