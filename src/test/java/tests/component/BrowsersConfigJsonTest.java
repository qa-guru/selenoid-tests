package tests.component;

import annotations.Layer;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertTrue;

@Layer("component")
@DisplayName("Browsers config fixture")
class BrowsersConfigJsonTest {

    @Test
    @DisplayName("parses browsers catalog map")
    void parsesBrowsersCatalog() {
        var json = FixtureJson.load("fixtures/ui/browsers-config.json");
        var versions = io.restassured.path.json.JsonPath.from(json).getMap("chrome");
        assertTrue(versions.containsKey("148.0"));
        assertTrue(versions.containsKey("149.0"));
    }
}
