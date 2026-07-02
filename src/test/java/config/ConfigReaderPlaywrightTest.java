package config;

import annotations.Component;
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

@Component("playwright-image")
@Layer("unit")
@DisplayName("ConfigReader Playwright endpoint")
@Execution(ExecutionMode.SAME_THREAD)
class ConfigReaderPlaywrightTest {

    private static TestConfig configWith(Map<String, String> overrides) {
        return ConfigFactory.create(TestConfig.class, overrides);
    }

    @Test
    @DisplayName("resolvePlaywrightWsEndpoint appends selenoid query params")
    void resolvePlaywrightWsEndpointAddsQueryParams() {
        var config = configWith(Map.of(
                "playwrightWsEndpoint", "ws://127.0.0.1:4444/playwright/playwright-chromium/1.61.1",
                "playwrightSessionName", "smoke",
                "playwrightSessionTimeout", "5m"
        ));
        var endpoint = ConfigReader.resolvePlaywrightWsEndpoint(config);
        assertTrue(endpoint.contains("name=smoke"));
        assertTrue(endpoint.contains("sessionTimeout=5m"));
        assertTrue(endpoint.contains("enableVNC=false"));
        assertTrue(endpoint.contains("enableVideo=false"));
    }

    @Test
    @DisplayName("resolvePlaywrightWsEndpoint fails fast when playwrightWsEndpoint is empty")
    void resolvePlaywrightWsEndpointFailsWhenEmpty() {
        var config = configWith(Map.of("playwrightWsEndpoint", ""));
        var error = assertThrows(IllegalStateException.class, () -> ConfigReader.resolvePlaywrightWsEndpoint(config));
        assertTrue(error.getMessage().contains("playwrightWsEndpoint"));
    }

    @Test
    @DisplayName("resolvePlaywrightWsEndpoint keeps endpoint that already has query string")
    void resolvePlaywrightWsEndpointKeepsExistingQueryString() {
        var preset = "ws://127.0.0.1:4444/playwright/playwright-chromium/1.61.1?name=preset&enableVNC=true";
        var config = configWith(Map.of(
                "playwrightWsEndpoint", preset,
                "playwrightSessionName", "ignored",
                "playwrightEnableVnc", "false"
        ));
        assertEquals(preset, ConfigReader.resolvePlaywrightWsEndpoint(config));
    }

    @Test
    @DisplayName("resolvePlaywrightWsEndpoint appends VNC and video flags from config")
    void resolvePlaywrightWsEndpointAppendsVncAndVideoFlags() {
        var config = configWith(Map.of(
                "playwrightWsEndpoint", "ws://127.0.0.1:4444/playwright/playwright-chromium/1.61.1",
                "playwrightSessionName", "rec",
                "playwrightSessionTimeout", "10m",
                "playwrightEnableVnc", "true",
                "playwrightEnableVideo", "true"
        ));
        var endpoint = ConfigReader.resolvePlaywrightWsEndpoint(config);
        assertTrue(endpoint.contains("enableVNC=true"));
        assertTrue(endpoint.contains("enableVideo=true"));
        assertTrue(endpoint.contains("sessionTimeout=10m"));
    }

    @Test
    @DisplayName("resolvePlaywrightWsEndpoint URL-encodes session name")
    void resolvePlaywrightWsEndpointUrlEncodesSessionName() {
        var config = configWith(Map.of(
                "playwrightWsEndpoint", "ws://127.0.0.1:4444/playwright/playwright-chromium/1.61.1",
                "playwrightSessionName", "my session"
        ));
        var endpoint = ConfigReader.resolvePlaywrightWsEndpoint(config);
        assertTrue(endpoint.contains("name=my+session"));
    }
}
