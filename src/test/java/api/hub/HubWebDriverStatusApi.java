package api.hub;

import config.ConfigReader;
import io.qameta.allure.Step;

import static io.restassured.RestAssured.given;

public final class HubWebDriverStatusApi {

    private HubWebDriverStatusApi() {
    }

    @Step("GET /wd/hub/status")
    public static HubWebDriverStatus fetch() {
        return fetchFrom(ConfigReader.resolveApiBaseUrl());
    }

    @Step("GET {baseUri}/wd/hub/status")
    public static HubWebDriverStatus fetchFrom(String baseUri) {
        return given()
                .baseUri(trimTrailingSlash(baseUri))
                .when()
                .get("/wd/hub/status")
                .then()
                .statusCode(200)
                .extract()
                .as(HubWebDriverStatus.class);
    }

    private static String trimTrailingSlash(String baseUri) {
        return baseUri.endsWith("/") ? baseUri.substring(0, baseUri.length() - 1) : baseUri;
    }
}
