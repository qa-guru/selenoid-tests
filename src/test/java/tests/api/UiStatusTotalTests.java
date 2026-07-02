package tests.api;

import annotations.Component;
import annotations.Layer;
import api.UiApiTestBase;
import api.ui.UiStatusApi;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import static io.qameta.allure.Allure.step;
import static org.junit.jupiter.api.Assertions.assertTrue;

@Layer("api")
@Component("selenoid-ui")
@Epic("selenoid-ui")
@Feature("UI status proxy")
@DisplayName("UI status total counter")
class UiStatusTotalTests extends UiApiTestBase {

    @Test
    @Tag("api")
    @Tag("positive")
    @DisplayName("GET UI /status exposes total quota counter")
    void uiStatusExposesTotalCounter() {
        var status = step("GET UI /status", UiStatusApi::fetch);
        step("Verify total counter is non-negative", () -> assertTrue(status.total() >= 0));
    }
}
