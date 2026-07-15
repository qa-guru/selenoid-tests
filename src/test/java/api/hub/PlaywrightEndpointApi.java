package api.hub;

import config.ConfigReader;
import io.qameta.allure.Step;

import java.net.URI;

import static io.restassured.RestAssured.given;

public final class PlaywrightEndpointApi {

    private PlaywrightEndpointApi() {
    }

    @Step("GET {path} without WebSocket upgrade")
    public static void assertUpgradeRequired() {
        var path = resolveHttpPath();
        hubRequest()
                .when()
                .get(path)
                .then()
                .statusCode(400);
    }

    static String resolveHttpPath() {
        var endpoint = ConfigReader.testConfig.playwrightWsEndpoint().trim();
        var httpUri = endpoint.replaceFirst("^ws:", "http:").replaceFirst("^wss:", "https:");
        var uri = URI.create(httpUri);
        var path = uri.getPath();
        if (path == null || path.isBlank()) {
            throw new IllegalStateException("Playwright path not found in: " + endpoint);
        }
        return withRawQuery(path, uri.getRawQuery());
    }

    @Step("GET unknown playwright path without WebSocket upgrade")
    public static void assertUnknownPathRejected() {
        hubRequest()
                .when()
                .get(withRawQuery("/playwright/unknown-browser/0.0.0", playwrightRawQuery()))
                .then()
                .statusCode(400);
    }

    private static String playwrightRawQuery() {
        var endpoint = ConfigReader.testConfig.playwrightWsEndpoint().trim();
        var httpUri = endpoint.replaceFirst("^ws:", "http:").replaceFirst("^wss:", "https:");
        return URI.create(httpUri).getRawQuery();
    }

    private static String withRawQuery(String path, String rawQuery) {
        if (rawQuery == null || rawQuery.isBlank()) {
            return path;
        }
        return path + "?" + rawQuery;
    }

    private static io.restassured.specification.RequestSpecification hubRequest() {
        var base = ConfigReader.resolveApiBaseUrl();
        var trimmed = base.endsWith("/") ? base.substring(0, base.length() - 1) : base;
        return given().baseUri(trimmed);
    }
}
