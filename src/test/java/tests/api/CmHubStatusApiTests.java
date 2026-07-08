package tests.api;

import annotations.Component;
import annotations.Layer;
import api.CmApiTestBase;
import api.hub.HubStatusApi;
import config.ConfigReader;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import static io.qameta.allure.Allure.step;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

@Layer("api")
@Component("cm")
@Epic("cm")
@Feature("CM-managed hub")
@DisplayName("CM hub status API")
class CmHubStatusApiTests extends CmApiTestBase {

    @Test
    @Tag("api")
    @Tag("cm")
    @Tag("positive")
    @DisplayName("GET /status on CM hub port returns statistics JSON")
    void cmHubStatusReturnsStatisticsJson() {
        var status = step("GET CM hub /status", () ->
                HubStatusApi.fetchFrom(ConfigReader.resolveCmHubUrl()));

        step("Verify hub counters are non-negative", () -> {
            assertTrue(status.total() >= 0);
            assertTrue(status.used() >= 0);
            assertNotNull(status.browsers());
        });
    }
}
