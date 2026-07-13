package tests.api;

import annotations.Component;
import annotations.Layer;
import api.ApiTestBase;
import api.hub.PlaywrightEndpointApi;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import io.qameta.allure.Story;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import static io.qameta.allure.Allure.step;

@Layer("api")
@Component("playwright-image")
@Epic("playwright-image")
@Feature("Playwright endpoint")
@Story("Playwright unknown path")
@DisplayName("Playwright unknown path")
class PlaywrightUnknownPathTests extends ApiTestBase {

    @Test
    @Tag("api")
    @Tag("negative")
    @DisplayName("GET unknown playwright path returns 400")
    void unknownPlaywrightPathRejected() {
        step("GET unknown playwright path", PlaywrightEndpointApi::assertUnknownPathRejected);
    }
}
