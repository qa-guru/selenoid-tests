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

@Component("webdriver-image")
@Layer("unit")
@DisplayName("ConfigReader WebDriver browser URL")
@Execution(ExecutionMode.SAME_THREAD)
class ConfigReaderWebdriverTest {

    private static TestConfig configWith(Map<String, String> overrides) {
        return ConfigFactory.create(TestConfig.class, overrides);
    }

    @Test
    @DisplayName("resolveUiBrowserUrl strips trailing slash when remoteUrl is blank")
    void resolveUiBrowserUrlUsesUiUrlWhenRemoteUrlBlank() {
        var config = configWith(Map.of(
                "uiUrl", "http://127.0.0.1:8080/",
                "remoteUrl", ""
        ));
        assertEquals("http://127.0.0.1:8080", ConfigReader.resolveUiBrowserUrl(config));
    }

    @Test
    @DisplayName("resolveUiBrowserUrl maps loopback to host.docker.internal for remote sessions")
    void resolveUiBrowserUrlMapsLoopbackForRemoteSessions() {
        var config = configWith(Map.of(
                "uiUrl", "http://127.0.0.1:8080",
                "remoteUrl", "http://127.0.0.1:4444/wd/hub"
        ));
        var browserUrl = ConfigReader.resolveUiBrowserUrl(config);
        assertTrue(browserUrl.contains("host.docker.internal"));
        assertTrue(browserUrl.endsWith(":8080"));
    }

    @Test
    @DisplayName("resolveUiBrowserUrl maps localhost for remote sessions")
    void resolveUiBrowserUrlMapsLocalhostForRemoteSessions() {
        var config = configWith(Map.of(
                "uiUrl", "http://localhost:8080/",
                "remoteUrl", "http://127.0.0.1:4444/wd/hub"
        ));
        assertEquals("http://host.docker.internal:8080", ConfigReader.resolveUiBrowserUrl(config));
    }

    @Test
    @DisplayName("resolveUiBrowserUrl fails fast when uiUrl is empty with remoteUrl set")
    void resolveUiBrowserUrlFailsWhenUiUrlEmptyWithRemoteUrl() {
        var config = configWith(Map.of(
                "uiUrl", "",
                "remoteUrl", "http://127.0.0.1:4444/wd/hub"
        ));
        var error = assertThrows(IllegalStateException.class, () -> ConfigReader.resolveUiBrowserUrl(config));
        assertTrue(error.getMessage().contains("uiUrl"));
    }
}
