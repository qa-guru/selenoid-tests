package api.ui;

import config.ConfigReader;
import io.qameta.allure.Step;

import java.util.Map;

import static io.restassured.RestAssured.given;

public final class UiBrowsersConfigApi {

    private UiBrowsersConfigApi() {
    }

    @Step("GET /browsers-config")
    public static Map<String, Object> fetch() {
        return fetchFrom(ConfigReader.resolveUiUrl());
    }

    @Step("GET {uiUrl}browsers-config")
    public static Map<String, Object> fetchFrom(String uiUrl) {
        var base = uiUrl.endsWith("/") ? uiUrl.substring(0, uiUrl.length() - 1) : uiUrl;
        return given()
                .baseUri(base)
                .when()
                .get("/browsers-config")
                .then()
                .statusCode(200)
                .extract()
                .as(Map.class);
    }
}
