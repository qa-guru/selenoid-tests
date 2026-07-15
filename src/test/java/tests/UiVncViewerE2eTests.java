package tests;

import annotations.Component;
import annotations.Layer;
import api.hub.HubSessionApi;
import config.ConfigReader;
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

import static io.qameta.allure.Allure.step;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertTrue;

@Layer("e2e")
@Component("selenoid-ui")
@Epic("selenoid-ui")
@Feature("VNC viewer")
@Story("VNC viewer in UI")
@DisplayName("UI VNC viewer e2e")
@ResourceLock(value = "hubSessions", mode = ResourceAccessMode.READ_WRITE)
class UiVncViewerE2eTests extends UiTestBase {

    private final UiCapabilitiesPage uiCapabilities = new UiCapabilitiesPage();
    private final UiSessionPage uiSession = new UiSessionPage();

    @Test
    @Tag("smoke")
    @Tag("positive")
    @DisplayName("Capabilities chrome session shows connected VNC and navigates remote browser")
    void vncViewerHappyPath() {
        var remoteUiUrl = ConfigReader.resolveUiUrl();
        String sessionId = null;
        try {
            step("Open dashboard and wait for CONNECTED", () ->
                    uiDashboard.openPage().shouldBeConnected());

            step("Create chrome session from Capabilities", () ->
                    uiCapabilities.openPage()
                            .selectChrome(config.chromeVersion())
                            .createSession());

            step("Wait for VNC connected and unlock screen", () ->
                    uiSession.shouldVncBeConnected().unlockVncScreen());

            sessionId = step("Read session id from URL", uiSession::sessionId);
            var hubSessionId = sessionId;

            step("Navigate remote browser via hub WebDriver API", () ->
                    HubSessionApi.navigate(hubSessionId, remoteUiUrl));

            step("Verify remote page loaded", () -> {
                var expectedHost = java.net.URI.create(remoteUiUrl).getHost();
                assertTrue(
                        HubSessionApi.getCurrentUrl(hubSessionId).contains(expectedHost),
                        "Remote browser URL should point to UI host");
                assertFalse(
                        HubSessionApi.getTitle(hubSessionId).isBlank(),
                        "Remote browser title should not be empty");
            });
        } finally {
            if (sessionId != null) {
                var finalSessionId = sessionId;
                step("Delete hub session", () -> HubSessionApi.delete(finalSessionId));
            }
        }
    }
}
