package tests.api;

import annotations.Component;
import annotations.Layer;
import api.UiApiTestBase;
import api.ui.UiProxyApi;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import static io.qameta.allure.Allure.step;

@Layer("api")
@Component("selenoid-ui")
@Epic("selenoid-ui")
@Feature("UI clipboard proxy")
@DisplayName("UI /clipboard proxy API")
class UiClipboardApiTests extends UiApiTestBase {

    @Test
    @Tag("api")
    @Tag("negative")
    @DisplayName("GET /clipboard/{sessionId} via UI proxy returns 404 for unknown session")
    void uiClipboardUnknownSessionReturnsNotFound() {
        step("GET UI /clipboard/unknown-session", () ->
                UiProxyApi.getExpectStatus("/clipboard", "unknown-session", 404));
    }
}
