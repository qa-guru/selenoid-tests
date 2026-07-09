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
        var alwaysMatch = new LinkedHashMap<String, Object>();
        alwaysMatch.put("browserName", browserName);
        alwaysMatch.put("browserVersion", browserVersion);
        alwaysMatch.put("goog:chromeOptions", Map.of("args", dockerChromeArgs()));
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
        var alwaysMatch = new LinkedHashMap<String, Object>();
        alwaysMatch.put("browserName", browserName);
        alwaysMatch.put("browserVersion", browserVersion);
        alwaysMatch.put("goog:chromeOptions", Map.of("args", dockerChromeArgs()));

        return Map.of("capabilities", Map.of("alwaysMatch", alwaysMatch));
    }

    private static List<String> dockerChromeArgs() {
        return ChromeOptions.dockerHeadless().args();
    }

    private static RequestSpecification hubRequest() {
        return given().baseUri(HubRequest.baseUri());
    }
}
