package tests.api;

import annotations.Component;
import annotations.Layer;
import api.ApiTestBase;
import api.hub.HubWelcomeApi;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import static io.qameta.allure.Allure.step;
import static org.junit.jupiter.api.Assertions.assertTrue;

@Layer("api")
@Component("selenoid")
@Epic("selenoid")
@Feature("Hub welcome")
@DisplayName("Hub welcome page API")
class HubWelcomeApiTests extends ApiTestBase {

    @Test
    @Tag("api")
    @Tag("positive")
    @DisplayName("GET / returns Selenoid welcome text")
    void welcomeReturnsSelenoidBanner() {
        var body = step("GET /", HubWelcomeApi::fetchText);
        step("Verify welcome text", () -> assertTrue(body.contains("Selenoid")));
    }
}
