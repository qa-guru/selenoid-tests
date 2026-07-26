package api.hub;

import io.qameta.allure.Step;

import java.time.Duration;
import java.util.List;
import java.util.concurrent.TimeUnit;

import static io.restassured.RestAssured.given;

/** Hub `/har/` artifact API — same listing shape as `/video/` (`videos` JSON field). */
public final class HubHarApi {

    private HubHarApi() {
    }

    @Step("GET /har/?json (paginated)")
    public static VideoListResponse listJson() {
        return listJson(10, 0, null);
    }

    @Step("GET /har/?json&limit={limit}&offset={offset}&q={q}")
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
                .get("/har/")
                .then()
                .statusCode(200)
                .extract()
                .as(VideoListResponse.class);
    }

    @Step("GET /har/?json&q={sessionId} — find session HAR name")
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

    @Step("GET /har/{fileName}")
    public static byte[] download(String fileName) {
        return given()
                .baseUri(HubRequest.baseUri())
                .when()
                .get("/har/{fileName}", fileName)
                .then()
                .statusCode(200)
                .extract()
                .asByteArray();
    }

    @Step("Wait for /har artifact for session {sessionId}")
    public static String waitForSessionHar(String sessionId, Duration timeout) throws InterruptedException {
        var deadline = System.currentTimeMillis() + timeout.toMillis();
        while (System.currentTimeMillis() < deadline) {
            var match = findBySessionId(sessionId);
            if (match != null) {
                return match;
            }
            TimeUnit.MILLISECONDS.sleep(400);
        }
        return findBySessionId(sessionId);
    }
}
