package tests;

import annotations.Layer;
import config.ConfigReader;
import helpers.UiClient;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import static io.qameta.allure.Allure.step;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

@Layer("integration")
@Epic("Selenoid")
@Feature("UI SSE")
@DisplayName("UI SSE stream")
class UiSseStreamTests {

    @Test
    @Tag("api")
    @Tag("positive")
    @DisplayName("SSE /events delivers state or error payload")
    void sseDeliversPayload() {
        var uiUrl = ConfigReader.resolveUiUrl();

        var json = step("Read first SSE event", () -> UiClient.readFirstSseEvent(uiUrl));

        step("Verify payload has state or errors", () -> {
            assertTrue(json.has("state") || json.has("errors"),
                    () -> "Expected state or errors in SSE payload: " + json);
            if (json.has("state")) {
                var state = json.getAsJsonObject("state");
                assertNotNull(state.get("total"));
                assertNotNull(state.get("browsers"));
            }
        });
    }
}
