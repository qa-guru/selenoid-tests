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

@Component("webdriver-image")
@Layer("component")
@Epic("webdriver-image")
@DisplayName("Edge min browser catalog fixture")
class MsedgeMinCatalogJsonTest {

    @Test
    @Tag("min")
    @DisplayName("parses msedge-min catalog version block port and path")
    void parsesMsedgeMinCatalogVersionBlock() {
        var json = FixtureJson.load("fixtures/webdriver/browser-catalog.json");
        var path = io.restassured.path.json.JsonPath.from(json);
        @SuppressWarnings("unchecked")
        var versionBlock = (Map<String, Object>) path.getMap("msedge.versions").get("145.0-min");
        assertEquals("4444", versionBlock.get("port"));
        assertEquals("/", versionBlock.get("path"));
    }

    @Test
    @Tag("min")
    @DisplayName("parses min image tag in msedge catalog entry")
    void parsesMinImageTagInMsedgeCatalog() {
        var json = FixtureJson.load("fixtures/webdriver/browser-catalog.json");
        var path = io.restassured.path.json.JsonPath.from(json);
        @SuppressWarnings("unchecked")
        var versionBlock = (Map<String, Object>) path.getMap("msedge.versions").get("145.0-min");
        assertTrue(String.valueOf(versionBlock.get("image")).contains("145-min"));
    }
}
