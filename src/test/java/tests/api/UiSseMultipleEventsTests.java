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
import static org.junit.jupiter.api.Assertions.assertTrue;

@Layer("api")
@Component("selenoid-ui")
@Epic("selenoid-ui")
@Feature("UI SSE")
@Story("UI SSE multiple events")
@DisplayName("UI SSE multiple events")
class UiSseMultipleEventsTests extends UiApiTestBase {

    @Test
    @Tag("api")
    @Tag("positive")
    @DisplayName("SSE /events delivers two consecutive payloads")
    void sseDeliversTwoEvents() {
        var events = step("Read two SSE events", () -> SseStreamApi.readTwoEvents());
        step("Verify both events are parseable", () -> {
            assertTrue(events.get(0).hasState() || events.get(0).hasErrors());
            assertTrue(events.get(1).hasState() || events.get(1).hasErrors());
        });
    }
}
