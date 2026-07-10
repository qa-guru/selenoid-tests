package tests.component;

import annotations.Component;
import annotations.Layer;
import io.qameta.allure.Epic;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertTrue;

@Component("cm")
@Layer("component")
@Epic("cm")
@DisplayName("CM browsers.json fixture")
class CmBrowsersConfigJsonTest {

    private static final String CI_BROWSERS = "fixtures/ci-browsers.json";

    @Test
    @DisplayName("parses cm browsers.json default chrome version")
    void parsesCmBrowsersDefaultVersion() {
        var json = FixtureJson.loadProject(CI_BROWSERS);
        var path = io.restassured.path.json.JsonPath.from(json);
        var defaultVersion = path.getString("chrome.default");
        assertTrue(defaultVersion.matches("\\d+\\.\\d+"));

        var root = path.getMap("");
        assertTrue(root.containsKey("firefox"));
        assertTrue(root.containsKey("msedge"));

        var versions = path.getMap("chrome.versions");
        assertTrue(versions.containsKey(defaultVersion));
        assertTrue(versions.containsKey(defaultVersion + "-min"));
    }

    @Test
    @DisplayName("parses cm browsers.json chrome-min image tag")
    void parsesCmBrowsersChromeMinImage() {
        var json = FixtureJson.loadProject(CI_BROWSERS);
        var path = io.restassured.path.json.JsonPath.from(json);
        var defaultVersion = path.getString("chrome.default");
        var minVersion = defaultVersion + "-min";
        var major = defaultVersion.substring(0, defaultVersion.indexOf('.'));
        var image = path.getString("chrome.versions.'%s'.image".formatted(minVersion));
        assertTrue(image.contains("webdriver-chrome"));
        assertTrue(image.contains(major + "-min"));
    }
}
