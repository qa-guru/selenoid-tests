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
        given()
                .when()
                .get(path)
                .then()
                .statusCode(400);
    }

    static String resolveHttpPath() {
        var endpoint = ConfigReader.testConfig.playwrightWsEndpoint().trim();
        var httpUri = endpoint.replaceFirst("^ws:", "http:").replaceFirst("^wss:", "https:");
        var path = URI.create(httpUri).getPath();
        if (path == null || path.isBlank()) {
            throw new IllegalStateException("Playwright path not found in: " + endpoint);
        }
        return path;
    }
}
