package tests.api;

import annotations.Component;
import annotations.Layer;
import api.UiApiTestBase;
import api.ui.UiPingApi;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import io.qameta.allure.Story;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import static io.qameta.allure.Allure.step;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

@Layer("api")
@Component("selenoid-ui")
@Epic("selenoid-ui")
@Feature("UI ping")
@Story("UI ping API")
@DisplayName("UI ping API")
class UiPingTests extends UiApiTestBase {

    @Test
    @Tag("api")
    @Tag("positive")
    @DisplayName("GET /ping returns uptime and version")
    void pingReturnsMetadata() {
        var ping = step("GET /ping", () -> UiPingApi.fetch());

        step("Verify uptime and version", () -> {
            assertNotNull(ping.uptime());
            assertFalse(ping.uptime().isBlank());
            assertNotNull(ping.version());
            assertFalse(ping.version().isBlank());
        });
    }
}
