package api.ui;

import io.qameta.allure.Step;

import java.util.List;

import static io.restassured.RestAssured.given;

public final class UiVideoApi {

    private UiVideoApi() {
    }

    @Step("GET UI /video/?json")
    public static List<String> listJson() {
        return given()
                .baseUri(UiRequest.baseUri())
                .queryParam("json", "")
                .when()
                .get("/video/")
                .then()
                .statusCode(200)
                .extract()
                .jsonPath()
                .getList("", String.class);
    }

    @Step("GET UI /video/{fileName} — expect HTTP {expectedStatus}")
    public static void getExpectStatus(String fileName, int expectedStatus) {
        given()
                .baseUri(UiRequest.baseUri())
                .when()
                .get("/video/{fileName}", fileName)
                .then()
                .statusCode(expectedStatus);
    }
}
