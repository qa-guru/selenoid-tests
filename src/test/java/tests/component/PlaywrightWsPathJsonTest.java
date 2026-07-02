package tests.component;

import annotations.Component;
import annotations.Layer;
import io.qameta.allure.Epic;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import java.util.Map;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertTrue;

@Component("playwright-image")
@Layer("component")
@Epic("playwright-image")
@DisplayName("Playwright WS path fixture")
class PlaywrightWsPathJsonTest {

    @Test
    @DisplayName("parses playwright browser catalog default version")
    void parsesPlaywrightDefaultVersion() {
        var json = FixtureJson.load("fixtures/playwright/browser-catalog.json");
        var defaultVersion = io.restassured.path.json.JsonPath.from(json)
                .getString("playwright-chromium.default");
        assertEquals("1.61.1", defaultVersion);
    }

    @Test
    @DisplayName("parses playwright protocol flag in catalog entry")
    void parsesPlaywrightProtocolFlag() {
        var json = FixtureJson.load("fixtures/playwright/browser-catalog.json");
        var path = io.restassured.path.json.JsonPath.from(json);
        @SuppressWarnings("unchecked")
        var versionBlock = (Map<String, Object>) path.getMap("playwright-chromium.versions").get("1.61.1");
        assertEquals("playwright", versionBlock.get("protocol"));
        assertTrue(String.valueOf(versionBlock.get("image")).contains("playwright-chromium"));
    }
}
