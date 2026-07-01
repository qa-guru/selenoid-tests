package api.hub;

import config.ConfigReader;
import io.qameta.allure.Step;

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
        return given()
                .baseUri(trimTrailingSlash(baseUri))
                .when()
                .get("/status")
                .then()
                .statusCode(200)
                .extract()
                .as(HubStatus.class);
    }

    private static String trimTrailingSlash(String baseUri) {
        return baseUri.endsWith("/") ? baseUri.substring(0, baseUri.length() - 1) : baseUri;
    }
}
