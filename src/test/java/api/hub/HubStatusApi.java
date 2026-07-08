package api.hub;

import config.ConfigReader;
import io.qameta.allure.Step;
import io.restassured.response.ResponseBodyExtractionOptions;

import static io.restassured.RestAssured.given;

public final class HubStatusApi {

    private HubStatusApi() {
    }

    @Step("GET /status")
    public static HubStatus fetch() {
        return fetchFrom(ConfigReader.resolveApiBaseUrl());
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

    private static String trimTrailingSlash(String baseUri) {
        return baseUri.endsWith("/") ? baseUri.substring(0, baseUri.length() - 1) : baseUri;
    }
}
