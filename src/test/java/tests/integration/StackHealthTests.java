package tests.integration;

import annotations.Component;
import annotations.Layer;
import api.hub.HubStatusApi;
import api.ui.UiPingApi;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import io.qameta.allure.Story;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import static io.qameta.allure.Allure.step;
import static org.junit.jupiter.api.Assertions.assertFalse;

@Layer("integration")
@Component("selenoid-ui")
@Epic("selenoid-ui")
@Feature("Stack health")
@Story("Stack health integration")
@DisplayName("Stack health integration")
class StackHealthTests {

    @Test
    @Tag("integration")
    @Tag("positive")
    @DisplayName("Hub and UI respond on health endpoints")
    void hubAndUiRespondOnHealthEndpoints() {
        step("GET hub /status", () -> HubStatusApi.fetch());
        var ping = step("GET UI /ping", UiPingApi::fetch);
        step("Verify UI ping uptime", () -> assertFalse(ping.uptime().isBlank()));
    }
}
