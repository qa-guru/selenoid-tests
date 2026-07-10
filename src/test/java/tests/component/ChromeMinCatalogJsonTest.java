package tests.component;

import annotations.Component;
import annotations.Layer;
import config.WebDriverCatalog;
import io.qameta.allure.Epic;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import java.util.Map;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertTrue;

@Component("webdriver-image")
@Layer("component")
@Epic("webdriver-image")
@DisplayName("Chrome min browser catalog fixture")
class ChromeMinCatalogJsonTest {

    @Test
    @Tag("min")
    @DisplayName("parses chrome-min catalog version block port and path")
    void parsesChromeMinCatalogVersionBlock() {
        var minVersion = WebDriverCatalog.minVersion("chrome");
        var json = FixtureJson.load(WebDriverCatalog.CATALOG_RESOURCE);
        var path = io.restassured.path.json.JsonPath.from(json);
        @SuppressWarnings("unchecked")
        var versionBlock = (Map<String, Object>) path.getMap("chrome.versions").get(minVersion);
        assertEquals("4444", versionBlock.get("port"));
        assertEquals("/", versionBlock.get("path"));
    }

    @Test
    @Tag("min")
    @DisplayName("parses min image tag in chrome catalog entry")
    void parsesMinImageTagInChromeCatalog() {
        var minVersion = WebDriverCatalog.minVersion("chrome");
        var json = FixtureJson.load(WebDriverCatalog.CATALOG_RESOURCE);
        var path = io.restassured.path.json.JsonPath.from(json);
        @SuppressWarnings("unchecked")
        var versionBlock = (Map<String, Object>) path.getMap("chrome.versions").get(minVersion);
        assertTrue(String.valueOf(versionBlock.get("image")).contains(WebDriverCatalog.minImageMajor("chrome") + "-min"));
    }
}
