package tests.api;

import annotations.Component;
import annotations.Layer;
import api.ApiTestBase;
import api.hub.HubWebDriverStatusApi;
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
@Feature("Hub WebDriver status")
@DisplayName("Hub WebDriver status API")
class HubWebDriverStatusApiTests extends ApiTestBase {

    @Test
    @Tag("api")
    @Tag("positive")
    @DisplayName("GET /wd/hub/status reports ready hub")
    void webDriverStatusReportsReady() {
        var status = step("GET /wd/hub/status", HubWebDriverStatusApi::fetch);

        step("Verify ready flag and message", () -> {
            assertNotNull(status.value());
            assertTrue(status.value().ready(), "Hub should be ready for new sessions");
            assertNotNull(status.value().message());
        });
    }
}
