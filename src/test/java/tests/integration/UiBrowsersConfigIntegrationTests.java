package tests.integration;

import annotations.Component;
import annotations.Layer;
import api.hub.HubStatusApi;
import api.ui.UiBrowsersConfigApi;
import config.ConfigReader;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import io.qameta.allure.Story;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import java.util.Map;

import static io.qameta.allure.Allure.step;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

@Layer("integration")
@Component("selenoid-ui")
@Epic("selenoid-ui")
@Feature("UI browsers config proxy")
@Story("UI browsers config vs hub status")
@DisplayName("UI browsers config vs hub status")
class UiBrowsersConfigIntegrationTests {

    @Test
    @Tag("integration")
    @Tag("positive")
    @DisplayName("UI /browsers-config includes chrome versions from hub /status")
    void browsersConfigIncludesHubBrowserVersions() {
        var hubBrowsers = step("GET hub /status browsers", () -> HubStatusApi.fetch().browsers());
        var uiCatalog = step("GET UI /browsers-config", UiBrowsersConfigApi::fetch);

        step("Verify configured chrome version appears in hub status and UI catalog", () -> {
            assertNotNull(hubBrowsers.get("chrome"), "Hub should expose chrome family");
            assertTrue(uiCatalog.containsKey("chrome"), "UI catalog should expose chrome family");
            var version = ConfigReader.testConfig.browserVersion();
            @SuppressWarnings("unchecked")
            var hubChrome = (Map<String, Object>) hubBrowsers.get("chrome");
            @SuppressWarnings("unchecked")
            var uiChrome = (Map<String, Object>) uiCatalog.get("chrome");
            assertTrue(hubChrome.containsKey(version),
                    () -> "Hub status missing chrome version " + version);
            assertTrue(uiChrome.containsKey(version),
                    () -> "UI catalog missing chrome version " + version);
        });
    }
}
