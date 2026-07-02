package tests.api;

import annotations.Component;
import annotations.Layer;
import api.ApiTestBase;
import api.hub.HubStatusApi;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import static io.qameta.allure.Allure.step;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

@Layer("api")
@Component("selenoid")
@Epic("selenoid")
@Feature("Hub status")
@DisplayName("Hub status")
class HubStatusTests extends ApiTestBase {

    @Test
    @Tag("api")
    @Tag("positive")
    @DisplayName("GET /status returns hub statistics JSON")
    void statusReturnsHubStatisticsJson() {
        var status = step("GET /status", () -> HubStatusApi.fetch());

        step("Verify root counters", () -> {
            assertTrue(status.total() >= 0);
            assertTrue(status.used() >= 0);
            assertTrue(status.queued() >= 0);
            assertTrue(status.pending() >= 0);
        });

        step("Verify browsers map is present", () -> assertNotNull(status.browsers()));
    }
}
