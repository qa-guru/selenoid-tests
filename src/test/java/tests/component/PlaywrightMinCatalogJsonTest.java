package tests.component;

import annotations.Component;
import annotations.Layer;
import io.qameta.allure.Epic;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import java.util.Map;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertTrue;

@Component("playwright-image")
@Layer("component")
@Epic("playwright-image")
@DisplayName("Playwright min browser catalog fixture")
class PlaywrightMinCatalogJsonTest {

    @Test
    @Tag("min")
    @DisplayName("parses playwrightVersion from min catalog version block")
    void parsesPlaywrightVersionFromMinCatalog() {
        var json = FixtureJson.load("fixtures/playwright/browser-catalog.json");
        var path = io.restassured.path.json.JsonPath.from(json);
        @SuppressWarnings("unchecked")
        var versionBlock = (Map<String, Object>) path.getMap("playwright-chromium.versions").get("1.61.1-min");
        assertEquals("1.61.1", versionBlock.get("playwrightVersion"));
    }

    @Test
    @Tag("min")
    @DisplayName("parses min image tag in catalog entry")
    void parsesMinImageTagInCatalog() {
        var json = FixtureJson.load("fixtures/playwright/browser-catalog.json");
        var path = io.restassured.path.json.JsonPath.from(json);
        @SuppressWarnings("unchecked")
        var versionBlock = (Map<String, Object>) path.getMap("playwright-chromium.versions").get("1.61.1-min");
        assertEquals("playwright", versionBlock.get("protocol"));
        assertTrue(String.valueOf(versionBlock.get("image")).contains("1.61.1-min"));
    }
}
