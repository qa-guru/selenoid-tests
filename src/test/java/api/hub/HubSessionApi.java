package api.hub;

import config.ConfigReader;
import config.TestConfig;
import io.qameta.allure.Step;
import io.restassured.http.ContentType;

import static io.restassured.RestAssured.given;

public final class HubSessionApi {

    private HubSessionApi() {
    }

    @Step("POST /wd/hub/session ({browserName} {browserVersion})")
    public static String create(String browserName, String browserVersion) {
        var request = new CreateSessionRequest(
                new CreateSessionCapabilities(new SessionAlwaysMatch(browserName, browserVersion)));
        return given()
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
        given()
                .when()
                .delete("/wd/hub/session/{sessionId}", sessionId)
                .then()
                .statusCode(200);
    }
}
