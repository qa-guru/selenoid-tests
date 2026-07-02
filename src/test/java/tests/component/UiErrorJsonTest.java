package tests.component;

import annotations.Layer;
import api.ui.UiErrorResponse;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertNotNull;

@Layer("component")
@DisplayName("UI error fixture")
class UiErrorJsonTest {

    @Test
    @DisplayName("parses error list payload")
    void parsesErrorPayload() {
        var error = FixtureJson.parse("fixtures/ui/error.json", UiErrorResponse.class);
        assertNotNull(error.errors());
        assertFalse(error.errors().isEmpty());
    }
}
