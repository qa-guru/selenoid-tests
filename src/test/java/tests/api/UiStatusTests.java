package tests.api;

import annotations.Layer;
import api.UiApiTestBase;
import api.ui.UiStatusApi;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import static io.qameta.allure.Allure.step;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

@Layer("api")
@Epic("selenoid-ui")
@Feature("UI status proxy")
@DisplayName("UI status proxy")
class UiStatusTests extends UiApiTestBase {

    @Test
    @Tag("api")
    @Tag("positive")
    @DisplayName("GET UI /status returns proxied hub statistics JSON")
    void uiStatusReturnsProxiedHubJson() {
        var status = step("GET UI /status", () -> UiStatusApi.fetch());

        step("Verify proxied state counters", () -> {
            assertTrue(status.total() >= 0);
            assertTrue(status.used() >= 0);
            assertTrue(status.queued() >= 0);
            assertTrue(status.pending() >= 0);
        });

        step("Verify browsers map is present", () -> {
            assertNotNull(status.browsers());
        });
    }
}
