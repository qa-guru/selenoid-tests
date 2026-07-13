package tests.integration;

import annotations.Component;
import annotations.Layer;
import api.hub.HubSessionApi;
import api.hub.HubStatusApi;
import api.ui.UiStatusApi;
import config.ConfigReader;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import io.qameta.allure.Story;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.parallel.ResourceAccessMode;
import org.junit.jupiter.api.parallel.ResourceLock;

import static io.qameta.allure.Allure.step;
import static org.junit.jupiter.api.Assertions.assertEquals;

@Layer("integration")
@Component("selenoid-ui")
@Epic("selenoid-ui")
@Feature("UI hub proxy")
@Story("UI hub proxy")
@DisplayName("UI status with active session")
@ResourceLock(value = "hubSessions", mode = ResourceAccessMode.READ_WRITE)
class UiStatusWithSessionTests {

    @Test
    @Tag("integration")
    @Tag("positive")
    @DisplayName("UI /status used counter increments while hub session is active")
    void uiUsedCounterIncrementsWithHubSession() {
        var usedBefore = step("Snapshot hub used counter", () -> HubStatusApi.fetch().used());
        var sessionId = step("Create hub session", () -> HubSessionApi.create());
        try {
            var hubUsed = step("GET hub /status used", () -> HubStatusApi.fetch().used());
            var uiUsed = step("GET UI /status used", () ->
                    UiStatusApi.fetchFrom(ConfigReader.resolveUiUrl()).used());
            step("Verify hub and UI report incremented used counter", () -> {
                assertEquals(usedBefore + 1, hubUsed);
                assertEquals(hubUsed, uiUsed);
            });
        } finally {
            step("Delete hub session", () -> HubSessionApi.delete(sessionId));
        }
        step("Verify used counter returned to baseline", () ->
                assertEquals(usedBefore, HubStatusApi.fetch().used()));
    }
}
