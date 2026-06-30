package tests;

import annotations.Layer;
import config.ConfigReader;
import helpers.HubClient;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import static io.qameta.allure.Allure.step;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

@Layer("integration")
@Epic("Selenoid")
@Feature("UI status proxy")
@DisplayName("UI status proxy")
class UiStatusTests {

    @Test
    @Tag("api")
    @Tag("positive")
    @DisplayName("GET UI /status returns proxied hub statistics JSON")
    void uiStatusReturnsProxiedHubJson() {
        var uiUrl = ConfigReader.resolveUiUrl();

        var response = step("GET UI /status", () -> HubClient.fetchStatus(uiUrl));

        step("Verify HTTP 200", () -> assertEquals(200, response.httpStatus()));

        var json = step("Parse JSON body", response::json);

        step("Verify proxied state counters", () -> {
            var state = json.has("state") ? json.getAsJsonObject("state") : json;
            assertCounter(state, "total");
            assertCounter(state, "used");
            assertCounter(state, "queued");
            assertCounter(state, "pending");
        });

        step("Verify browsers map is present", () -> {
            var state = json.has("state") ? json.getAsJsonObject("state") : json;
            assertNotNull(state.get("browsers"));
            assertTrue(state.get("browsers").isJsonObject());
        });
    }

    private static void assertCounter(com.google.gson.JsonObject json, String field) {
        assertNotNull(json.get(field), () -> "Missing field: " + field);
        assertTrue(json.get(field).isJsonPrimitive(), () -> "Field is not numeric: " + field);
    }
}
