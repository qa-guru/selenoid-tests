package tests.api;

import annotations.Layer;
import api.ApiTestBase;
import api.hub.HubSessionApi;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import static io.qameta.allure.Allure.step;
import static org.junit.jupiter.api.Assertions.assertFalse;

@Layer("api")
@Epic("selenoid")
@Feature("WebDriver session API")
@DisplayName("Hub session API")
class HubSessionApiTests extends ApiTestBase {

    @Test
    @Tag("api")
    @Tag("positive")
    @DisplayName("POST /wd/hub/session creates and DELETE removes session")
    void createAndDeleteSession() {
        var sessionId = step("Create remote session", () -> HubSessionApi.create(config));

        step("Verify session id", () -> assertFalse(sessionId.isBlank()));

        step("Delete session", () -> HubSessionApi.delete(sessionId));
    }
}
