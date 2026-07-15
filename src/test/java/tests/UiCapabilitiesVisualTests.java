package tests;

import annotations.Component;
import annotations.Layer;
import helpers.ScreenshotBaseline;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import io.qameta.allure.Story;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;
import pages.UiCapabilitiesPage;

import static com.codeborne.selenide.Selenide.$;
import static io.qameta.allure.Allure.step;

@Layer("e2e")
@Component("selenoid-ui")
@Epic("selenoid-ui")
@Feature("Capabilities visual")
@Story("Capabilities select and Create Session preserve Selenoid 2 look")
@DisplayName("UI Capabilities visual")
class UiCapabilitiesVisualTests extends UiTestBase {

    private final UiCapabilitiesPage uiCapabilities = new UiCapabilitiesPage();

    @Test
    @Tag("visual")
    @DisplayName("Closed browser select and disabled Create Session match baseline")
    void closedSelectAndDisabledCreateSessionMatchBaseline() {
        step("Open Capabilities", () ->
                uiCapabilities.openPage());

        step("Compare closed select + disabled Create Session", () ->
                ScreenshotBaseline.captureAndCompare(
                        uiCapabilities.setupPanel(),
                        "capabilities-closed",
                        viewportWidth(),
                        "capabilities-closed-select-button"));
    }

    @Test
    @Tag("visual")
    @DisplayName("Open browser select menu matches baseline")
    void openBrowserSelectMenuMatchesBaseline() {
        step("Open Capabilities and expand browser menu", () ->
                uiCapabilities.openPage().openBrowserMenu());

        step("Compare open select menu", () ->
                ScreenshotBaseline.captureAndCompare(
                        $(".Select__menu"),
                        "capabilities-open",
                        viewportWidth(),
                        "capabilities-open-select-menu"));
    }

    @Test
    @Tag("visual")
    @DisplayName("Enabled Create Session after browser select matches baseline")
    void enabledCreateSessionAfterSelectMatchesBaseline() {
        step("Open Capabilities and select chrome", () ->
                uiCapabilities.openPage().selectChrome(config.chromeVersion()));

        step("Compare enabled Create Session", () ->
                ScreenshotBaseline.captureAndCompare(
                        uiCapabilities.setupPanel(),
                        "capabilities-selected",
                        viewportWidth(),
                        "capabilities-selected-create-session"));
    }

    private int viewportWidth() {
        var size = config.browserSize();
        var separator = size.indexOf('x');
        if (separator <= 0) {
            return 1920;
        }
        return Integer.parseInt(size.substring(0, separator));
    }
}
