package api.hub;

import config.ConfigReader;
import io.qameta.allure.Step;

import java.util.List;

import static io.restassured.RestAssured.given;

public final class HubVideoApi {

    private HubVideoApi() {
    }

    @Step("GET /video/?json")
    public static List<String> listJson() {
        return given()
                .baseUri(HubRequest.baseUri())
                .queryParam("json", "")
                .when()
                .get("/video/")
                .then()
                .statusCode(200)
                .extract()
                .jsonPath()
                .getList("", String.class);
    }

    @Step("GET /video/{fileName}")
    public static byte[] download(String fileName) {
        return given()
                .baseUri(HubRequest.baseUri())
                .when()
                .get("/video/{fileName}", fileName)
                .then()
                .statusCode(200)
                .extract()
                .asByteArray();
    }

    @Step("GET /video/{fileName}")
    public static void getExpectStatus(String fileName, int expectedStatus) {
        given()
                .baseUri(HubRequest.baseUri())
                .when()
                .get("/video/{fileName}", fileName)
                .then()
                .statusCode(expectedStatus);
    }
}
