package tests.api;

import annotations.Component;
import annotations.Layer;
import api.ApiTestBase;
import api.hub.HubErrorApi;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import io.qameta.allure.Story;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import static io.qameta.allure.Allure.step;

@Layer("api")
@Component("selenoid")
@Epic("selenoid")
@Feature("Hub error")
@Story("Hub /error API")
@DisplayName("Hub /error API")
class HubErrorApiTests extends ApiTestBase {

    @Test
    @Tag("api")
    @Tag("negative")
    @DisplayName("GET /error returns invalid session id JSON")
    void errorReturnsInvalidSessionJson() {
        step("GET /error", HubErrorApi::fetchExpectInvalidSessionJson);
    }
}
