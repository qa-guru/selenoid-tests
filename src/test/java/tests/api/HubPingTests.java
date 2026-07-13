package tests.api;

import annotations.Component;
import annotations.Layer;
import api.ApiTestBase;
import api.hub.HubPingApi;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import io.qameta.allure.Story;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import static io.qameta.allure.Allure.step;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertNotNull;

@Layer("api")
@Component("selenoid")
@Epic("selenoid")
@Feature("Hub ping")
@Story("Hub ping API")
@DisplayName("Hub ping API")
class HubPingTests extends ApiTestBase {

    @Test
    @Tag("api")
    @Tag("positive")
    @DisplayName("GET /ping returns uptime and version")
    void pingReturnsMetadata() {
        var ping = step("GET /ping", HubPingApi::fetch);

        step("Verify uptime and version", () -> {
            assertNotNull(ping.uptime());
            assertFalse(ping.uptime().isBlank());
            assertNotNull(ping.version());
            assertFalse(ping.version().isBlank());
        });
    }
}
