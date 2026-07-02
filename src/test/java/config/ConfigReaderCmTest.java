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

@Component("cm")
@Layer("unit")
@DisplayName("ConfigReader CM URLs")
@Execution(ExecutionMode.SAME_THREAD)
class ConfigReaderCmTest {

    private static TestConfig configWith(Map<String, String> overrides) {
        return ConfigFactory.create(TestConfig.class, overrides);
    }

    @Test
    @DisplayName("resolveCmHubUrl uses cmHubPort with trailing slash")
    void resolveCmHubUrlUsesCmHubPort() {
        var config = configWith(Map.of("cmHubPort", "4445"));
        assertEquals("http://127.0.0.1:4445/", ConfigReader.resolveCmHubUrl(config));
    }

    @Test
    @DisplayName("resolveCmHubUrl respects custom cmHubPort")
    void resolveCmHubUrlRespectsCustomPort() {
        var config = configWith(Map.of("cmHubPort", "9444"));
        assertEquals("http://127.0.0.1:9444/", ConfigReader.resolveCmHubUrl(config));
    }

    @Test
    @DisplayName("resolveCmUiUrl uses cmUiPort with trailing slash")
    void resolveCmUiUrlUsesCmUiPort() {
        var config = configWith(Map.of("cmUiPort", "8081"));
        assertEquals("http://127.0.0.1:8081/", ConfigReader.resolveCmUiUrl(config));
    }

    @Test
    @DisplayName("resolveCmUiUrl respects custom cmUiPort")
    void resolveCmUiUrlRespectsCustomPort() {
        var config = configWith(Map.of("cmUiPort", "9080"));
        assertEquals("http://127.0.0.1:9080/", ConfigReader.resolveCmUiUrl(config));
    }

    @Test
    @DisplayName("resolveCmRemoteUrl points WebDriver hub at CM-managed port")
    void resolveCmRemoteUrlCombinesHubBaseAndWdHubPath() {
        var config = configWith(Map.of("cmHubPort", "4445"));
        assertEquals("http://127.0.0.1:4445/wd/hub", ConfigReader.resolveCmRemoteUrl(config));
    }
}
