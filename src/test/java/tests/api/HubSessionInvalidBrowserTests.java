package tests.api;

import annotations.Layer;
import api.ApiTestBase;
import api.hub.HubSessionApi;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import static io.qameta.allure.Allure.step;

@Layer("api")
@Epic("selenoid")
@Feature("WebDriver session API")
@DisplayName("Hub session invalid browser")
class HubSessionInvalidBrowserTests extends ApiTestBase {

    @Test
    @Tag("api")
    @Tag("negative")
    @DisplayName("POST /wd/hub/session rejects unknown browser")
    void createSessionRejectsUnknownBrowser() {
        step("POST session with unknown browser", () ->
                HubSessionApi.createExpectStatus("unknown-browser", "1.0", 500));
    }
}
