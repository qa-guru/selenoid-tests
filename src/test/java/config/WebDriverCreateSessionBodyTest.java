package config;

import annotations.Component;
import annotations.Layer;
import api.hub.HubSessionApi;
import com.google.gson.Gson;
import io.restassured.path.json.JsonPath;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.util.Map;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertTrue;

@Component("webdriver-image")
@Layer("unit")
@DisplayName("WebDriver createSession body")
class WebDriverCreateSessionBodyTest {

    private static final Gson GSON = new Gson();
    private static final String CATALOG = "fixtures/webdriver/browser-catalog.json";

    @Test
    @DisplayName("firefox warm session body matches catalog default version")
    void firefoxWarmSessionBodyMatchesCatalog() {
        var warmVersion = catalogDefaultVersion("firefox");
        var body = HubSessionApi.createSessionBody("firefox", warmVersion);
        var json = GSON.toJson(body);

        assertSessionBody(json, "firefox", warmVersion, "moz:firefoxOptions", "-headless");
    }

    @Test
    @DisplayName("msedge warm session body matches catalog default version")
    void msedgeWarmSessionBodyMatchesCatalog() {
        var warmVersion = catalogDefaultVersion("msedge");
        var body = HubSessionApi.createSessionBody("msedge", warmVersion);
        var json = GSON.toJson(body);

        assertSessionBody(json, "msedge", warmVersion, "ms:edgeOptions", "no-sandbox");
    }

    @Test
    @DisplayName("chrome warm session body matches catalog default version")
    void chromeWarmSessionBodyMatchesCatalog() {
        var warmVersion = catalogDefaultVersion("chrome");
        var body = HubSessionApi.createSessionBody("chrome", warmVersion);
        var json = GSON.toJson(body);

        assertSessionBody(json, "chrome", warmVersion, "goog:chromeOptions", "no-sandbox");
    }

    @Test
    @DisplayName("chrome min session body includes docker-safe chrome args")
    void chromeMinSessionBodyIncludesDockerSafeChromeArgs() {
        var minVersion = "149.0-min";
        var body = HubSessionApi.createSessionBody("chrome", minVersion);
        var json = GSON.toJson(body);

        assertSessionBody(json, "chrome", minVersion, "goog:chromeOptions", "no-sandbox");
        assertTrue(catalogVersionBlock("chrome", minVersion).get("image").toString().contains("149-min"));
    }

    @Test
    @DisplayName("firefox min session body includes docker-safe firefox args")
    void firefoxMinSessionBodyIncludesDockerSafeFirefoxArgs() {
        var minVersion = "151.0-min";
        var body = HubSessionApi.createSessionBody("firefox", minVersion);
        var json = GSON.toJson(body);

        assertSessionBody(json, "firefox", minVersion, "moz:firefoxOptions", "-headless");
        assertTrue(catalogVersionBlock("firefox", minVersion).get("image").toString().contains("151-min"));
    }

    @Test
    @DisplayName("msedge min session body includes docker-safe edge args")
    void msedgeMinSessionBodyIncludesDockerSafeEdgeArgs() {
        var minVersion = "145.0-min";
        var body = HubSessionApi.createSessionBody("msedge", minVersion);
        var json = GSON.toJson(body);

        assertSessionBody(json, "msedge", minVersion, "ms:edgeOptions", "no-sandbox");
        assertTrue(catalogVersionBlock("msedge", minVersion).get("image").toString().contains("145-min"));
    }

    private static String catalogDefaultVersion(String browser) {
        return JsonPath.from(loadCatalog()).getString(browser + ".default");
    }

    @SuppressWarnings("unchecked")
    private static Map<String, Object> catalogVersionBlock(String browser, String version) {
        return (Map<String, Object>) JsonPath.from(loadCatalog()).getMap(browser + ".versions").get(version);
    }

    private static String loadCatalog() {
        try (var stream = WebDriverCreateSessionBodyTest.class.getClassLoader().getResourceAsStream(CATALOG)) {
            if (stream == null) {
                throw new IllegalStateException("Fixture not found: " + CATALOG);
            }
            return new String(stream.readAllBytes(), StandardCharsets.UTF_8);
        } catch (IOException e) {
            throw new IllegalStateException("Failed to read fixture: " + CATALOG, e);
        }
    }

    private static void assertSessionBody(String json, String browserName, String browserVersion,
            String optionsKey, String optionsArgFragment) {
        assertTrue(json.contains("capabilities"));
        assertTrue(json.contains("alwaysMatch"));
        assertTrue(json.contains("browserName"));
        assertTrue(json.contains(browserName));
        assertTrue(json.contains("browserVersion"));
        assertTrue(json.contains(browserVersion));
        assertTrue(json.contains(optionsKey));
        assertTrue(json.contains(optionsArgFragment));
        assertEquals(browserName, JsonPath.from(json).getString("capabilities.alwaysMatch.browserName"));
        assertEquals(browserVersion, JsonPath.from(json).getString("capabilities.alwaysMatch.browserVersion"));
    }
}
