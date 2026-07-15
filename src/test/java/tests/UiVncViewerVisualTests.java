package tests;

import annotations.Component;
import annotations.Layer;
import api.hub.HubSessionApi;
import helpers.ScreenshotBaseline;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import io.qameta.allure.Story;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.parallel.ResourceAccessMode;
import org.junit.jupiter.api.parallel.ResourceLock;
import pages.UiCapabilitiesPage;
import pages.UiSessionPage;

import static com.codeborne.selenide.Selenide.$;
import static io.qameta.allure.Allure.step;

@Layer("e2e")
@Component("selenoid-ui")
@Epic("selenoid-ui")
@Feature("VNC viewer")
@Story("VNC viewer visual")
@DisplayName("UI VNC viewer visual")
@ResourceLock(value = "hubSessions", mode = ResourceAccessMode.READ_WRITE)
class UiVncViewerVisualTests extends UiTestBase {

    private static final String BASELINE_AREA = "vnc-session";

    private final UiCapabilitiesPage uiCapabilities = new UiCapabilitiesPage();
    private final UiSessionPage uiSession = new UiSessionPage();

    @Test
    @Tag("visual")
    @DisplayName("VNC-connected session screen matches baseline")
    void vncConnectedSessionScreenMatchesBaseline() {
        String sessionId = null;
        try {
            step("Open dashboard and wait for CONNECTED", () ->
                    uiDashboard.openPage().shouldBeConnected());

            step("Create chrome session from Capabilities", () ->
                    uiCapabilities.openPage()
                            .selectChrome(config.chromeVersion())
                            .createSession());

            step("Wait for VNC connected", uiSession::shouldVncBeConnected);

            sessionId = step("Read session id from URL", uiSession::sessionId);

            step("Compare VNC card screenshot with baseline", () ->
                    ScreenshotBaseline.captureAndCompare(
                            $(".vnc-card"),
                            BASELINE_AREA,
                            viewportWidth(),
                            "vnc-connected-session"));
        } finally {
            if (sessionId != null) {
                var finalSessionId = sessionId;
                step("Delete hub session", () -> HubSessionApi.delete(finalSessionId));
            }
        }
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
