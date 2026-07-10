package config;

import annotations.Component;
import annotations.Layer;
import api.hub.HubSessionApi;
import com.google.gson.Gson;
import io.restassured.path.json.JsonPath;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import java.util.Map;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertTrue;

@Component("webdriver-image")
@Layer("unit")
@DisplayName("WebDriver createSession body")
class WebDriverCreateSessionBodyTest {

    private static final Gson GSON = new Gson();

    @Test
    @DisplayName("firefox warm session body matches catalog default version")
    void firefoxWarmSessionBodyMatchesCatalog() {
        var warmVersion = WebDriverCatalog.defaultVersion("firefox");
        var body = HubSessionApi.createSessionBody("firefox", warmVersion);
        var json = GSON.toJson(body);

        assertSessionBody(json, "firefox", warmVersion, "moz:firefoxOptions", "-headless");
    }

    @Test
    @DisplayName("msedge warm session body matches catalog default version")
    void msedgeWarmSessionBodyMatchesCatalog() {
        var warmVersion = WebDriverCatalog.defaultVersion("msedge");
        var body = HubSessionApi.createSessionBody("msedge", warmVersion);
        var json = GSON.toJson(body);

        assertSessionBody(json, "msedge", warmVersion, "ms:edgeOptions", "no-sandbox");
    }

    @Test
    @DisplayName("chrome warm session body matches catalog default version")
    void chromeWarmSessionBodyMatchesCatalog() {
        var warmVersion = WebDriverCatalog.defaultVersion("chrome");
        var body = HubSessionApi.createSessionBody("chrome", warmVersion);
        var json = GSON.toJson(body);

        assertSessionBody(json, "chrome", warmVersion, "goog:chromeOptions", "no-sandbox");
    }

    @Test
    @DisplayName("chrome min session body includes docker-safe chrome args")
    void chromeMinSessionBodyIncludesDockerSafeChromeArgs() {
        var minVersion = WebDriverCatalog.minVersion("chrome");
        var body = HubSessionApi.createSessionBody("chrome", minVersion);
        var json = GSON.toJson(body);

        assertSessionBody(json, "chrome", minVersion, "goog:chromeOptions", "no-sandbox");
        var image = WebDriverCatalog.versionBlock("chrome", minVersion).get("image").toString();
        assertTrue(image.contains(WebDriverCatalog.minImageMajor("chrome") + "-min"));
    }

    @Test
    @DisplayName("firefox min session body includes docker-safe firefox args")
    void firefoxMinSessionBodyIncludesDockerSafeFirefoxArgs() {
        var minVersion = WebDriverCatalog.minVersion("firefox");
        var body = HubSessionApi.createSessionBody("firefox", minVersion);
        var json = GSON.toJson(body);

        assertSessionBody(json, "firefox", minVersion, "moz:firefoxOptions", "-headless");
        var image = WebDriverCatalog.versionBlock("firefox", minVersion).get("image").toString();
        assertTrue(image.contains(WebDriverCatalog.minImageMajor("firefox") + "-min"));
    }

    @Test
    @DisplayName("msedge min session body includes docker-safe edge args")
    void msedgeMinSessionBodyIncludesDockerSafeEdgeArgs() {
        var minVersion = WebDriverCatalog.minVersion("msedge");
        var body = HubSessionApi.createSessionBody("msedge", minVersion);
        var json = GSON.toJson(body);

        assertSessionBody(json, "msedge", minVersion, "ms:edgeOptions", "no-sandbox");
        var image = WebDriverCatalog.versionBlock("msedge", minVersion).get("image").toString();
        assertTrue(image.contains(WebDriverCatalog.minImageMajor("msedge") + "-min"));
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
