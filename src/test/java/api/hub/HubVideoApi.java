package api.hub;

import io.qameta.allure.Step;

import java.util.List;

import static io.restassured.RestAssured.given;

public final class HubVideoApi {

    private HubVideoApi() {
    }

    @Step("GET /video/?json (paginated)")
    public static VideoListResponse listJson() {
        return listJson(10, 0, null);
    }

    @Step("GET /video/?json&limit={limit}&offset={offset}&q={q}")
    public static VideoListResponse listJson(int limit, int offset, String q) {
        var request = given()
                .baseUri(HubRequest.baseUri())
                .queryParam("json", "")
                .queryParam("limit", limit)
                .queryParam("offset", offset);
        if (q != null && !q.isBlank()) {
            request = request.queryParam("q", q);
        }
        return request
                .when()
                .get("/video/")
                .then()
                .statusCode(200)
                .extract()
                .as(VideoListResponse.class);
    }

    @Step("GET /video/?json — video names only")
    public static List<String> listNames() {
        return listJson().videos();
    }

    @Step("GET /video/?json&q={sessionId} — find session video name")
    public static String findBySessionId(String sessionId) {
        var listed = listJson(10, 0, sessionId);
        if (listed.videos() == null) {
            return null;
        }
        return listed.videos().stream()
                .filter(name -> name.contains(sessionId))
                .findFirst()
                .orElse(null);
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
