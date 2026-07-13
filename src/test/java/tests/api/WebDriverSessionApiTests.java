package tests.api;

import annotations.Component;
import annotations.Layer;
import api.ApiTestBase;
import api.hub.HubSessionApi;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import io.qameta.allure.Story;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.parallel.ResourceAccessMode;
import org.junit.jupiter.api.parallel.ResourceLock;

import static io.qameta.allure.Allure.step;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;

@Layer("api")
@Component("webdriver-image")
@Epic("webdriver-image")
@Feature("WebDriver session API")
@Story("WebDriver session API")
@DisplayName("WebDriver image session API")
@ResourceLock(value = "hubSessions", mode = ResourceAccessMode.READ_WRITE)
class WebDriverSessionApiTests extends ApiTestBase {

    @Test
    @Tag("api")
    @Tag("positive")
    @DisplayName("POST /wd/hub/session creates chrome session via webdriver-image node")
    void createAndDeleteSession() {
        var created = step("Create remote session", () -> HubSessionApi.createWithCapabilities(config));

        step("Verify session id and browserName", () -> {
            assertFalse(created.sessionId().isBlank());
            assertEquals(config.browser(), created.browserName());
        });

        step("Delete session", () -> HubSessionApi.delete(created.sessionId()));
    }
}
