package api.hub;

import config.ConfigReader;
import io.qameta.allure.Step;

import static io.restassured.RestAssured.given;

public final class HubWelcomeApi {

    private HubWelcomeApi() {
    }

    @Step("GET /")
    public static String fetchText() {
        return given()
                .baseUri(HubRequest.baseUri())
                .when()
                .get("/")
                .then()
                .statusCode(200)
                .extract()
                .asString();
    }
}
