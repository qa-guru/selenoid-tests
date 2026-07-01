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
@Epic("selenoid")
@Feature("Hub status")
@DisplayName("Hub status")
class HubStatusTests {

    @Test
    @Tag("api")
    @Tag("positive")
    @DisplayName("GET /status returns hub statistics JSON")
    void statusReturnsHubStatisticsJson() {
        var hubUrl = ConfigReader.resolveHubUrl();

        var response = step("GET /status", () -> HubClient.fetchStatus(hubUrl));

        step("Verify HTTP 200", () -> assertEquals(200, response.httpStatus()));

        var json = step("Parse JSON body", response::json);

        step("Verify root counters", () -> {
            assertCounter(json, "total");
            assertCounter(json, "used");
            assertCounter(json, "queued");
            assertCounter(json, "pending");
        });

        step("Verify browsers map is present", () -> {
            assertNotNull(json.get("browsers"));
            assertTrue(json.get("browsers").isJsonObject());
        });
    }

    private static void assertCounter(com.google.gson.JsonObject json, String field) {
        assertNotNull(json.get(field), () -> "Missing field: " + field);
        assertTrue(json.get(field).isJsonPrimitive(), () -> "Field is not numeric: " + field);
    }
}
