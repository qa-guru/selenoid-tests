package api.hub;

import config.ConfigReader;
import io.qameta.allure.Step;

import static io.restassured.RestAssured.given;

public final class HubErrorApi {

    private HubErrorApi() {
    }

    @Step("GET /error")
    public static void fetchExpectInvalidSessionJson() {
        given()
                .baseUri(HubRequest.baseUri())
                .when()
                .get("/error")
                .then()
                .statusCode(404)
                .body("value.error", org.hamcrest.Matchers.equalTo("invalid session id"));
    }
}
