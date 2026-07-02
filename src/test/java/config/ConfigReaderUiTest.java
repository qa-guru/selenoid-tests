package config;

import annotations.Layer;
import org.aeonbits.owner.ConfigFactory;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.parallel.Execution;
import org.junit.jupiter.api.parallel.ExecutionMode;

import java.util.Map;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

@Layer("unit")
@DisplayName("ConfigReader UI URL")
@Execution(ExecutionMode.SAME_THREAD)
class ConfigReaderUiTest {

    private static TestConfig configWith(Map<String, String> overrides) {
        return ConfigFactory.create(TestConfig.class, overrides);
    }

    @Test
    @DisplayName("resolveUiUrl fails fast when uiUrl is empty")
    void resolveUiUrlFailsWhenEmpty() {
        var config = configWith(Map.of("uiUrl", ""));
        var error = assertThrows(IllegalStateException.class, () -> ConfigReader.resolveUiUrl(config));
        assertTrue(error.getMessage().contains("uiUrl"));
    }

    @Test
    @DisplayName("resolveUiUrl keeps trailing slash")
    void resolveUiUrlKeepsTrailingSlash() {
        var config = configWith(Map.of("uiUrl", "http://127.0.0.1:8080/"));
        assertEquals("http://127.0.0.1:8080/", ConfigReader.resolveUiUrl(config));
    }
}
