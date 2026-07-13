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
@Feature("Hub clipboard proxy")
@Story("Hub clipboard proxy API")
@DisplayName("Hub clipboard proxy API")
class HubClipboardApiTests extends ApiTestBase {

    @Test
    @Tag("api")
    @Tag("negative")
    @DisplayName("GET /clipboard/{sessionId} for unknown session returns 404 JSON")
    void clipboardUnknownSessionReturnsNotFound() {
        step("GET /clipboard/unknown-session", () ->
                HubProxyApi.getExpectStatus("/clipboard", "unknown-session", 404));
    }
}
