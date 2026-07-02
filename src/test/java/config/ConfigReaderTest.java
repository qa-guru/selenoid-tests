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
@DisplayName("ConfigReader")
@Execution(ExecutionMode.SAME_THREAD)
class ConfigReaderTest {

    private static TestConfig configWith(Map<String, String> overrides) {
        return ConfigFactory.create(TestConfig.class, overrides);
    }

    @Test
    @DisplayName("resolveHubUrl adds trailing slash")
    void resolveHubUrlAddsTrailingSlash() {
        var config = configWith(Map.of("hubUrl", "http://127.0.0.1:4444"));
        assertEquals("http://127.0.0.1:4444/", ConfigReader.resolveHubUrl(config));
    }

    @Test
    @DisplayName("resolveHubUrl keeps trailing slash")
    void resolveHubUrlKeepsTrailingSlash() {
        var config = configWith(Map.of("hubUrl", "http://127.0.0.1:4444/"));
        assertEquals("http://127.0.0.1:4444/", ConfigReader.resolveHubUrl(config));
    }

    @Test
    @DisplayName("resolveHubUrl fails fast when hubUrl is empty")
    void resolveHubUrlFailsWhenEmpty() {
        var config = configWith(Map.of("hubUrl", ""));
        var error = assertThrows(IllegalStateException.class, () -> ConfigReader.resolveHubUrl(config));
        assertTrue(error.getMessage().contains("hubUrl"));
    }

    @Test
    @DisplayName("resolveUiUrl adds trailing slash")
    void resolveUiUrlAddsTrailingSlash() {
        var config = configWith(Map.of("uiUrl", "http://127.0.0.1:8080"));
        assertEquals("http://127.0.0.1:8080/", ConfigReader.resolveUiUrl(config));
    }

    @Test
    @DisplayName("resolveApiBaseUrl prefers apiBaseUrl over hubUrl")
    void resolveApiBaseUrlPrefersExplicitKey() {
        var config = configWith(Map.of(
                "apiBaseUrl", "http://api.example.com",
                "hubUrl", "http://hub.example.com/"));
        assertEquals("http://api.example.com/", ConfigReader.resolveApiBaseUrl(config));
    }

    @Test
    @DisplayName("resolveApiBaseUrl falls back to hubUrl when apiBaseUrl is empty")
    void resolveApiBaseUrlFallsBackToHubUrl() {
        var config = configWith(Map.of("apiBaseUrl", "", "hubUrl", "http://127.0.0.1:4444"));
        assertEquals("http://127.0.0.1:4444/", ConfigReader.resolveApiBaseUrl(config));
    }

    @Test
    @DisplayName("resolveApiBaseUrl fails fast when apiBaseUrl and hubUrl are empty")
    void resolveApiBaseUrlFailsWhenBothEmpty() {
        var config = configWith(Map.of("apiBaseUrl", "", "hubUrl", ""));
        var error = assertThrows(IllegalStateException.class, () -> ConfigReader.resolveApiBaseUrl(config));
        assertTrue(error.getMessage().contains("apiBaseUrl or hubUrl"));
    }
}
