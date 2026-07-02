package tests.api;

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

@Layer("api")
@Epic("selenoid")
@Feature("Hub status")
@DisplayName("Hub status browsers")
class HubStatusBrowsersTests extends ApiTestBase {

    @Test
    @Tag("api")
    @Tag("positive")
    @DisplayName("GET /status lists configured browser families")
    void statusListsBrowserFamilies() {
        var status = step("GET /status", HubStatusApi::fetch);
        step("Verify chrome is configured", () -> assertNotNull(status.browsers().get("chrome")));
    }
}
