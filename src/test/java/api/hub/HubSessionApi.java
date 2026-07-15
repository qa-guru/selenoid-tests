package api.hub;

import io.qameta.allure.Step;
import io.restassured.http.ContentType;
import io.restassured.specification.RequestSpecification;

import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

import config.ConfigReader;
import config.TestConfig;

import static io.restassured.RestAssured.given;

public final class HubSessionApi {

    private HubSessionApi() {
    }

    @Step("POST /wd/hub/session ({browserName} {browserVersion})")
    public static String create(String browserName, String browserVersion) {
        return hubRequest()
                .contentType(ContentType.JSON)
                .body(createSessionBody(browserName, browserVersion))
                .when()
                .post("/wd/hub/session")
                .then()
                .statusCode(200)
                .extract()
                .path("value.sessionId");
    }

    @Step("POST /wd/hub/session")
    public static String create(TestConfig config) {
        return create(config.browser(), config.browserVersion());
    }

    @Step("POST /wd/hub/session with selenoid:options")
    public static String createWithSelenoidOptions(TestConfig config, java.util.Map<String, Object> selenoidOptions) {
        return createWithSelenoidOptions(config.browser(), config.browserVersion(), selenoidOptions);
    }

    @Step("POST /wd/hub/session with selenoid:options ({browserVersion})")
    public static String createWithSelenoidOptions(String browserVersion, java.util.Map<String, Object> selenoidOptions) {
        return createWithSelenoidOptions("chrome", browserVersion, selenoidOptions);
    }

    @Step("POST /wd/hub/session with selenoid:options ({browserName} {browserVersion})")
    public static String createWithSelenoidOptions(String browserName, String browserVersion,
            java.util.Map<String, Object> selenoidOptions) {
        var alwaysMatch = createAlwaysMatch(browserName, browserVersion);
        if (selenoidOptions != null && !selenoidOptions.isEmpty()) {
            alwaysMatch.put("selenoid:options", selenoidOptions);
        }
        return hubRequest()
                .contentType(ContentType.JSON)
                .body(Map.of("capabilities", Map.of("alwaysMatch", alwaysMatch)))
                .when()
                .post("/wd/hub/session")
                .then()
                .statusCode(200)
                .extract()
                .path("value.sessionId");
    }

    @Step("POST /wd/hub/session")
    public static String create() {
        return create(ConfigReader.testConfig);
    }

    @Step("POST /wd/hub/session/{sessionId}/url — navigate to {url}")
    public static void navigate(String sessionId, String url) {
        hubRequest()
                .contentType(ContentType.JSON)
                .body(Map.of("url", url))
                .when()
                .post("/wd/hub/session/{sessionId}/url", sessionId)
                .then()
                .statusCode(200);
    }

    @Step("GET /wd/hub/session/{sessionId}/title")
    public static String getTitle(String sessionId) {
        return hubRequest()
                .when()
                .get("/wd/hub/session/{sessionId}/title", sessionId)
                .then()
                .statusCode(200)
                .extract()
                .path("value");
    }

    @Step("GET /wd/hub/session/{sessionId}/url")
    public static String getCurrentUrl(String sessionId) {
        return hubRequest()
                .when()
                .get("/wd/hub/session/{sessionId}/url", sessionId)
                .then()
                .statusCode(200)
                .extract()
                .path("value");
    }

    @Step("DELETE /wd/hub/session/{sessionId}")
    public static void delete(String sessionId) {
        hubRequest()
                .when()
                .delete("/wd/hub/session/{sessionId}", sessionId)
                .then()
                .statusCode(200);
    }

    @Step("DELETE /wd/hub/session/{sessionId} — expect HTTP {expectedStatus}")
    public static void deleteExpectStatus(String sessionId, int expectedStatus) {
        hubRequest()
                .when()
                .delete("/wd/hub/session/{sessionId}", sessionId)
                .then()
                .statusCode(expectedStatus);
    }

    @Step("POST /wd/hub/session — expect HTTP {expectedStatus}")
    public static void createExpectStatus(String browserName, String browserVersion, int expectedStatus) {
        hubRequest()
                .contentType(ContentType.JSON)
                .body(createSessionBody(browserName, browserVersion))
                .when()
                .post("/wd/hub/session")
                .then()
                .statusCode(expectedStatus);
    }

    @Step("POST /wd/hub/session and read browserName from capabilities")
    public static SessionCreateResult createWithCapabilities(TestConfig config) {
        var response = hubRequest()
                .contentType(ContentType.JSON)
                .body(createSessionBody(config.browser(), config.browserVersion()))
                .when()
                .post("/wd/hub/session")
                .then()
                .statusCode(200)
                .extract()
                .response();
        var sessionId = response.path("value.sessionId").toString();
        var browserName = response.path("value.capabilities.browserName");
        if (browserName == null) {
            browserName = response.path("value.capabilities.alwaysMatch.browserName");
        }
        return new SessionCreateResult(sessionId, browserName == null ? "" : browserName.toString());
    }

    public static Map<String, Object> createSessionBody(String browserName, String browserVersion) {
        return Map.of("capabilities", Map.of("alwaysMatch", createAlwaysMatch(browserName, browserVersion)));
    }

    private static LinkedHashMap<String, Object> createAlwaysMatch(String browserName, String browserVersion) {
        var alwaysMatch = new LinkedHashMap<String, Object>();
        alwaysMatch.put("browserName", browserName);
        alwaysMatch.put("browserVersion", browserVersion);
        switch (browserName) {
            case "firefox" -> alwaysMatch.put("moz:firefoxOptions", Map.of("args", dockerFirefoxArgs()));
            case "msedge" -> alwaysMatch.put("ms:edgeOptions", Map.of("args", dockerEdgeArgs()));
            default -> alwaysMatch.put("goog:chromeOptions", Map.of("args", dockerChromeArgs()));
        }
        return alwaysMatch;
    }

    private static List<String> dockerChromeArgs() {
        return ChromeOptions.dockerHeadless().args();
    }

    private static List<String> dockerFirefoxArgs() {
        return FirefoxOptions.dockerHeadless().args();
    }

    private static List<String> dockerEdgeArgs() {
        return EdgeOptions.dockerHeadless().args();
    }

    private static RequestSpecification hubRequest() {
        return given().baseUri(HubRequest.baseUri());
    }
}
