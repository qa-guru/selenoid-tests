package tests.component;

import annotations.Component;
import annotations.Layer;
import api.ui.UiStatusResponse;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;

@Component("selenoid-ui")
@Layer("component")
@DisplayName("UI status fixture")
class UiStatusJsonTest {

    @Test
    @DisplayName("parses proxied UI status wrapper")
    void parsesUiStatusWrapper() {
        var status = FixtureJson.parse("fixtures/ui/status.json", UiStatusResponse.class);
        assertNotNull(status.state());
        assertEquals(1, status.state().used());
    }
}
