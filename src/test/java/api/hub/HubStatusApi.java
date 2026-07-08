package api.hub;

import config.ConfigReader;
import io.qameta.allure.Step;

import java.time.Duration;

import static io.restassured.RestAssured.given;

public final class HubStatusApi {

    private HubStatusApi() {
    }

    @Step("GET /status")
    public static HubStatus fetch() {
        return fetchFrom(ConfigReader.resolveApiBaseUrl());
    }

    @Step("GET {baseUri}{hubStatusPath}")
    public static HubStatus fetchFrom(String baseUri) {
        return given()
                .baseUri(trimTrailingSlash(baseUri))
                .when()
                .get(ConfigReader.resolveHubStatusPath())
                .then()
                .statusCode(200)
                .extract()
                .body()
                .as(HubStatus.class);
    }

    @Step("Wait until hub used counter = {expectedUsed}")
    public static HubStatus waitUntilUsed(int expectedUsed, Duration timeout) {
        var deadline = System.currentTimeMillis() + timeout.toMillis();
        HubStatus last = null;
        while (System.currentTimeMillis() <= deadline) {
            last = fetch();
            if (last.used() == expectedUsed) {
                return last;
            }
            try {
                Thread.sleep(250);
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
                throw new IllegalStateException("Interrupted while waiting for used=" + expectedUsed, e);
            }
        }
        return last != null ? last : fetch();
    }

    private static String trimTrailingSlash(String baseUri) {
        return baseUri.endsWith("/") ? baseUri.substring(0, baseUri.length() - 1) : baseUri;
    }
}
