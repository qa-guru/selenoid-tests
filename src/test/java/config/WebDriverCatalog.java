package config;

import io.restassured.path.json.JsonPath;

import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.util.Map;

/**
 * SSOT for WebDriver browser versions in unit/component tests — reads {@code fixtures/webdriver/browser-catalog.json}.
 * Runtime integration/api/e2e use {@link TestConfig} (-D / config/*.properties).
 */
public final class WebDriverCatalog {

    public static final String CATALOG_RESOURCE = "fixtures/webdriver/browser-catalog.json";

    private WebDriverCatalog() {
    }

    public static String defaultVersion(String browser) {
        return JsonPath.from(loadCatalog()).getString(browser + ".default");
    }

    public static String minVersion(String browser) {
        return defaultVersion(browser) + "-min";
    }

    @SuppressWarnings("unchecked")
    public static Map<String, Object> versionBlock(String browser, String version) {
        return (Map<String, Object>) JsonPath.from(loadCatalog()).getMap(browser + ".versions").get(version);
    }

    public static String minImageMajor(String browser) {
        return defaultVersion(browser).substring(0, defaultVersion(browser).indexOf('.'));
    }

    static String loadCatalog() {
        try (var stream = WebDriverCatalog.class.getClassLoader().getResourceAsStream(CATALOG_RESOURCE)) {
            if (stream == null) {
                throw new IllegalStateException("Fixture not found: " + CATALOG_RESOURCE);
            }
            return new String(stream.readAllBytes(), StandardCharsets.UTF_8);
        } catch (IOException e) {
            throw new IllegalStateException("Failed to read fixture: " + CATALOG_RESOURCE, e);
        }
    }
}
