package api.hub;

import config.ConfigReader;
import io.qameta.allure.Step;

import static io.restassured.RestAssured.given;

public final class HubPingApi {

    private HubPingApi() {
    }

    @Step("GET /ping")
    public static HubPingResponse fetch() {
        return fetchFrom(ConfigReader.resolveApiBaseUrl());
    }

    @Step("GET {baseUri}/ping")
    public static HubPingResponse fetchFrom(String baseUri) {
        var base = baseUri.endsWith("/") ? baseUri.substring(0, baseUri.length() - 1) : baseUri;
        return given()
                .baseUri(base)
                .when()
                .get("/ping")
                .then()
                .statusCode(200)
                .extract()
                .as(HubPingResponse.class);
    }
}
