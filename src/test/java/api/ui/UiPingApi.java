package api.ui;

import io.qameta.allure.Step;

import static io.restassured.RestAssured.given;

public final class UiPingApi {

    private UiPingApi() {
    }

    @Step("GET /ping")
    public static UiPingResponse fetch() {
        return given()
                .when()
                .get("/ping")
                .then()
                .statusCode(200)
                .extract()
                .as(UiPingResponse.class);
    }
}
