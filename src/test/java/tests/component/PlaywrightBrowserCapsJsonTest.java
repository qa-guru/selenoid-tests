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
@DisplayName("Playwright browser caps fixture")
class PlaywrightBrowserCapsJsonTest {

    @Test
    @DisplayName("parses playwrightVersion from catalog version block")
    void parsesPlaywrightVersionFromCatalog() {
        var json = FixtureJson.load("fixtures/playwright/browser-catalog.json");
        var path = io.restassured.path.json.JsonPath.from(json);
        @SuppressWarnings("unchecked")
        var versionBlock = (Map<String, Object>) path.getMap("playwright-chromium.versions").get("1.61.1");
        assertEquals("1.61.1", versionBlock.get("playwrightVersion"));
    }

    @Test
    @DisplayName("lists playwright-chromium in browser families map")
    void listsPlaywrightChromiumFamily() {
        var json = FixtureJson.load("fixtures/playwright/browser-catalog.json");
        var families = io.restassured.path.json.JsonPath.from(json).getMap("$");
        assertTrue(families.containsKey("playwright-chromium"));
    }
}
