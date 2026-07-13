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
@Story("Playwright HTTP endpoint")
@DisplayName("Playwright HTTP endpoint")
class PlaywrightEndpointTests extends ApiTestBase {

    @Test
    @Tag("api")
    @Tag("positive")
    @DisplayName("GET playwright path without upgrade returns 400")
    void playwrightPathRequiresWebSocketUpgrade() {
        step("GET playwright path without WebSocket headers", PlaywrightEndpointApi::assertUpgradeRequired);
    }
}
