package tests.api;

import annotations.Component;
import annotations.Layer;
import api.CmApiTestBase;
import api.hub.HubSessionApi;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import io.qameta.allure.Story;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import static io.qameta.allure.Allure.step;
import static org.junit.jupiter.api.Assertions.assertFalse;

@Layer("api")
@Component("cm")
@Epic("cm")
@Feature("CM-managed hub session")
@Story("CM hub session API")
@DisplayName("CM hub session API")
class CmHubSessionApiTests extends CmApiTestBase {

    @Test
    @Tag("api")
    @Tag("cm")
    @Tag("positive")
    @DisplayName("POST /wd/hub/session on CM hub creates and DELETE removes session")
    void cmHubCreateAndDeleteSession() {
        var sessionId = step("Create session on CM hub", () -> HubSessionApi.create(config));

        step("Verify session id", () -> assertFalse(sessionId.isBlank()));

        step("Delete session on CM hub", () -> HubSessionApi.delete(sessionId));
    }
}
