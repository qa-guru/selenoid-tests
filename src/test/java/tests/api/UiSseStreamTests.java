package tests.api;

import annotations.Component;
import annotations.Layer;
import api.UiApiTestBase;
import api.ui.SseStreamApi;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import io.qameta.allure.Story;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import static io.qameta.allure.Allure.step;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

@Layer("api")
@Component("selenoid-ui")
@Epic("selenoid-ui")
@Feature("UI SSE")
@Story("UI SSE stream")
@DisplayName("UI SSE stream")
class UiSseStreamTests extends UiApiTestBase {

    @Test
    @Tag("api")
    @Tag("positive")
    @DisplayName("SSE /events delivers state or error payload")
    void sseDeliversPayload() {
        var event = step("Read first SSE event", () -> SseStreamApi.readFirstEvent());

        step("Verify payload has state or errors", () -> {
            assertTrue(event.hasState() || event.hasErrors(),
                    () -> "Expected state or errors in SSE payload: " + event);
            if (event.hasState()) {
                assertTrue(event.state().total() >= 0);
                assertNotNull(event.state().browsers());
            }
        });
    }
}
