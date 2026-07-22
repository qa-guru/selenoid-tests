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
@Story("UI SSE indicator")
@DisplayName("UI SSE indicator")
class UiSseIndicatorTests extends UiTestBase {

    @Test
    @Tag("smoke")
    @Tag("positive")
    @DisplayName("SSE indicator uses connected StatusTile styling when connected")
    void sseIndicatorUsesOkStyling() {
        step("Open dashboard", () -> uiDashboard.openPage().shouldBeConnected());
        step("Verify SSE indicator uses connected StatusTile styling", () ->
                $("#sse-status").shouldHave(com.codeborne.selenide.Condition.cssClass("status-tile--connected")));
    }
}
