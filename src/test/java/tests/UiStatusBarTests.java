package tests;

import annotations.Layer;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import static com.codeborne.selenide.Selenide.$;
import static io.qameta.allure.Allure.step;

@Layer("e2e")
@Epic("Selenoid")
@Feature("UI status bar")
@DisplayName("UI status bar")
class UiStatusBarTests extends UiTestBase {

    @Test
    @Tag("smoke")
    @Tag("positive")
    @DisplayName("SSE and SELENOID stay CONNECTED on stable stack")
    void statusBarStaysConnected() throws InterruptedException {
        step("Open dashboard and wait for CONNECTED", () ->
                uiDashboard.openPage().shouldBeConnected());

        step("Verify indicator styling", () -> {
            $("[data-testid='sse-status'] .indicator").shouldHave(com.codeborne.selenide.Condition.cssClass("indicator_ok"));
            $("[data-testid='selenoid-status'] .indicator").shouldHave(com.codeborne.selenide.Condition.cssClass("indicator_ok"));
        });

        step("Keep CONNECTED stable for 5 seconds", () ->
                uiDashboard.shouldStayConnected(5_000));
    }
}
