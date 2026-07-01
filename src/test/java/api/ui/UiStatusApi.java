package api.ui;

import api.hub.HubStatus;
import config.ConfigReader;
import io.qameta.allure.Step;
import io.restassured.response.ResponseBodyExtractionOptions;

import static io.restassured.RestAssured.given;

public final class UiStatusApi {

    private UiStatusApi() {
    }

    @Step("GET /status")
    public static HubStatus fetch() {
        return fetchFrom(ConfigReader.resolveUiUrl());
    }

    @Step("GET {baseUri}/status")
    public static HubStatus fetchFrom(String baseUri) {
        ResponseBodyExtractionOptions body = given()
                .baseUri(trimTrailingSlash(baseUri))
                .when()
                .get("/status")
                .then()
                .statusCode(200)
                .extract()
                .body();
        if (body.path("state") != null) {
            return body.jsonPath().getObject("state", HubStatus.class);
        }
        return body.as(HubStatus.class);
    }

    @Step("GET /status when hub is unavailable")
    public static UiErrorResponse fetchWhenHubUnavailable() {
        return given()
                .baseUri(trimTrailingSlash(ConfigReader.resolveUiUrl()))
                .when()
                .get("/status")
                .then()
                .statusCode(500)
                .extract()
                .as(UiErrorResponse.class);
    }

    private static String trimTrailingSlash(String baseUri) {
        return baseUri.endsWith("/") ? baseUri.substring(0, baseUri.length() - 1) : baseUri;
    }
}
