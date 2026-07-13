package tests.api;

import annotations.Component;
import annotations.Layer;
import api.ApiTestBase;
import api.hub.HubProxyApi;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import io.qameta.allure.Story;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import static io.qameta.allure.Allure.step;

@Layer("api")
@Component("selenoid")
@Epic("selenoid")
@Feature("Hub download proxy")
@Story("Hub download proxy API")
@DisplayName("Hub download proxy API")
class HubDownloadApiTests extends ApiTestBase {

    @Test
    @Tag("api")
    @Tag("negative")
    @DisplayName("GET /download/{sessionId} for unknown session returns 404 JSON")
    void downloadUnknownSessionReturnsNotFound() {
        step("GET /download/unknown-session", () ->
                HubProxyApi.getExpectStatus("/download", "unknown-session", 404));
    }
}
