package api.hub;

import config.ConfigReader;
import io.qameta.allure.Step;

import static io.restassured.RestAssured.given;

public final class HubProxyApi {

    private HubProxyApi() {
    }

    @Step("GET {prefix}/{sessionId} — expect HTTP {expectedStatus}")
    public static void getExpectStatus(String prefix, String sessionId, int expectedStatus) {
        given()
                .baseUri(HubRequest.baseUri())
                .when()
                .get(prefix + "/{sessionId}", sessionId)
                .then()
                .statusCode(expectedStatus);
    }
}
