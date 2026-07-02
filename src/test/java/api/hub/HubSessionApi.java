package api.hub;

import config.ConfigReader;
import config.TestConfig;
import io.qameta.allure.Step;
import io.restassured.http.ContentType;
import io.restassured.specification.RequestSpecification;

import static io.restassured.RestAssured.given;

public final class HubSessionApi {

    private HubSessionApi() {
    }

    @Step("POST /wd/hub/session ({browserName} {browserVersion})")
    public static String create(String browserName, String browserVersion) {
        var request = new CreateSessionRequest(
                new CreateSessionCapabilities(new SessionAlwaysMatch(browserName, browserVersion)));
        return hubRequest()
                .contentType(ContentType.JSON)
                .body(request)
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
        var request = new CreateSessionRequest(
                new CreateSessionCapabilities(new SessionAlwaysMatch(browserName, browserVersion)));
        hubRequest()
                .contentType(ContentType.JSON)
                .body(request)
                .when()
                .post("/wd/hub/session")
                .then()
                .statusCode(expectedStatus);
    }

    @Step("POST /wd/hub/session and read browserName from capabilities")
    public static SessionCreateResult createWithCapabilities(TestConfig config) {
        var request = new CreateSessionRequest(
                new CreateSessionCapabilities(new SessionAlwaysMatch(config.browser(), config.browserVersion())));
        var response = hubRequest()
                .contentType(ContentType.JSON)
                .body(request)
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

    private static RequestSpecification hubRequest() {
        return given().baseUri(trimTrailingSlash(ConfigReader.resolveApiBaseUrl()));
    }

    private static String trimTrailingSlash(String baseUri) {
        return baseUri.endsWith("/") ? baseUri.substring(0, baseUri.length() - 1) : baseUri;
    }
}
