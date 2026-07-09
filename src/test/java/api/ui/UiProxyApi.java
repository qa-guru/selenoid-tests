package api.ui;

import io.qameta.allure.Step;

import static io.restassured.RestAssured.given;

public final class UiProxyApi {

    private UiProxyApi() {
    }

    @Step("GET UI {prefix}/{sessionId} — expect HTTP {expectedStatus}")
    public static void getExpectStatus(String prefix, String sessionId, int expectedStatus) {
        given()
                .baseUri(UiRequest.baseUri())
                .when()
                .get(prefix + "/{sessionId}", sessionId)
                .then()
                .statusCode(expectedStatus);
    }
}
