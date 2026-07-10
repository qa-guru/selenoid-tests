package tests.component;

import annotations.Component;
import annotations.Layer;
import io.qameta.allure.Epic;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertTrue;

@Component("cm")
@Layer("component")
@Epic("cm")
@DisplayName("CM browsers.json fixture")
class CmBrowsersConfigJsonTest {

    @Test
    @DisplayName("parses cm browsers.json default chrome version")
    void parsesCmBrowsersDefaultVersion() {
        var json = FixtureJson.load("fixtures/cm/browsers-config.json");
        var defaultVersion = io.restassured.path.json.JsonPath.from(json).getString("chrome.default");
        assertEquals("149.0", defaultVersion);

        var root = io.restassured.path.json.JsonPath.from(json).getMap("");
        assertTrue(root.containsKey("firefox"));
        assertTrue(root.containsKey("msedge"));

        var versions = io.restassured.path.json.JsonPath.from(json).getMap("chrome.versions");
        assertTrue(versions.containsKey("149.0"));
        assertTrue(versions.containsKey("149.0-min"));
    }

    @Test
    @DisplayName("parses cm browsers.json chrome-min image tag")
    void parsesCmBrowsersChromeMinImage() {
        var json = FixtureJson.load("fixtures/cm/browsers-config.json");
        var image = io.restassured.path.json.JsonPath.from(json)
                .getString("chrome.versions.'149.0-min'.image");
        assertEquals("qaguru/webdriver-chrome:149-min", image);
    }
}
