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

import static io.qameta.allure.Allure.step;

@Layer("api")
@Component("selenoid")
@Epic("selenoid")
@Feature("WebDriver session API")
@Story("Hub session delete unknown")
@DisplayName("Hub session delete unknown")
class HubSessionDeleteUnknownTests extends ApiTestBase {

    @Test
    @Tag("api")
    @Tag("negative")
    @DisplayName("DELETE unknown session returns client error")
    void deleteUnknownSessionReturnsError() {
        step("DELETE unknown session id", () ->
                HubSessionApi.deleteExpectStatus("missing-session-id", 404));
    }
}
