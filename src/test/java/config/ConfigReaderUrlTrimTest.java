package config;

import annotations.Layer;
import org.aeonbits.owner.ConfigFactory;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.parallel.Execution;
import org.junit.jupiter.api.parallel.ExecutionMode;

import java.util.Map;

import static org.junit.jupiter.api.Assertions.assertEquals;

@Layer("unit")
@DisplayName("ConfigReader URL trim")
@Execution(ExecutionMode.SAME_THREAD)
class ConfigReaderUrlTrimTest {

    private static TestConfig configWith(Map<String, String> overrides) {
        return ConfigFactory.create(TestConfig.class, overrides);
    }

    @Test
    @DisplayName("resolveHubUrl trims surrounding whitespace")
    void resolveHubUrlTrimsWhitespace() {
        var config = configWith(Map.of("hubUrl", "  http://127.0.0.1:4444  "));
        assertEquals("http://127.0.0.1:4444/", ConfigReader.resolveHubUrl(config));
    }

    @Test
    @DisplayName("resolveUiUrl trims surrounding whitespace")
    void resolveUiUrlTrimsWhitespace() {
        var config = configWith(Map.of("uiUrl", "  http://127.0.0.1:8080  "));
        assertEquals("http://127.0.0.1:8080/", ConfigReader.resolveUiUrl(config));
    }

    @Test
    @DisplayName("resolveApiBaseUrl trims explicit apiBaseUrl")
    void resolveApiBaseUrlTrimsWhitespace() {
        var config = configWith(Map.of("apiBaseUrl", "  http://api.example.com  ", "hubUrl", "http://hub/"));
        assertEquals("http://api.example.com/", ConfigReader.resolveApiBaseUrl(config));
    }
}
