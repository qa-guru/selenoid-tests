package tests.integration;

import annotations.Layer;
import api.hub.HubStatusApi;
import api.ui.UiStatusApi;
import config.ConfigReader;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import static io.qameta.allure.Allure.step;
import static org.junit.jupiter.api.Assertions.assertEquals;

@Layer("integration")
@Epic("selenoid-ui")
@Feature("UI hub proxy")
@DisplayName("UI hub status consistency")
class UiHubStatusConsistencyTests {

    @Test
    @Tag("integration")
    @Tag("positive")
    @DisplayName("UI /status state mirrors hub /status counters")
    void uiStatusMirrorsHubCounters() {
        var hubStatus = step("GET hub /status", () ->
                HubStatusApi.fetchFrom(ConfigReader.resolveHubUrl()));
        var uiStatus = step("GET UI /status", () ->
                UiStatusApi.fetchFrom(ConfigReader.resolveUiUrl()));

        step("Verify proxied counters match hub", () -> {
            assertEquals(hubStatus.total(), uiStatus.total());
            assertEquals(hubStatus.used(), uiStatus.used());
            assertEquals(hubStatus.queued(), uiStatus.queued());
            assertEquals(hubStatus.pending(), uiStatus.pending());
        });
    }
}
