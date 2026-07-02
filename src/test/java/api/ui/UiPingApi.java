package api.ui;

import config.ConfigReader;
import io.qameta.allure.Step;

import static io.restassured.RestAssured.given;

public final class UiPingApi {

    private UiPingApi() {
    }

    @Step("GET /ping")
    public static UiPingResponse fetch() {
        return fetchFrom(ConfigReader.resolveUiUrl());
    }

    @Step("GET {baseUri}/ping")
    public static UiPingResponse fetchFrom(String baseUri) {
        var base = baseUri.endsWith("/") ? baseUri.substring(0, baseUri.length() - 1) : baseUri;
        return given()
                .baseUri(base)
                .when()
                .get("/ping")
                .then()
                .statusCode(200)
                .extract()
                .as(UiPingResponse.class);
    }
}
