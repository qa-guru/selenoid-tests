package tests.api;

import annotations.Component;
import annotations.Layer;
import api.ApiTestBase;
import api.hub.HubSessionApi;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.parallel.ResourceAccessMode;
import org.junit.jupiter.api.parallel.ResourceLock;

import static io.qameta.allure.Allure.step;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;

@Layer("api")
@Component("selenoid")
@Epic("selenoid")
@Feature("WebDriver session API")
@DisplayName("Hub session capabilities")
@ResourceLock(value = "hubSessions", mode = ResourceAccessMode.READ_WRITE)
class HubCapabilitiesApiTests extends ApiTestBase {

    @Test
    @Tag("api")
    @Tag("positive")
    @DisplayName("POST /wd/hub/session echoes browserName in capabilities")
    void createSessionEchoesBrowserName() {
        var created = step("Create remote session", () -> HubSessionApi.createWithCapabilities(config));

        step("Verify capabilities.browserName", () -> {
            assertFalse(created.browserName().isBlank());
            assertEquals(config.browser(), created.browserName());
        });

        step("Delete session", () -> HubSessionApi.delete(created.sessionId()));
    }
}
