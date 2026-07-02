package tests.api;

import annotations.Component;
import annotations.Layer;
import api.ApiTestBase;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import static io.qameta.allure.Allure.step;
import static io.restassured.RestAssured.given;
import static org.hamcrest.Matchers.containsString;

@Layer("api")
@Component("selenoid")
@Epic("selenoid")
@Feature("Hub logs")
@DisplayName("Hub logs endpoint API")
class HubLogsListApiTests extends ApiTestBase {

    @Test
    @Tag("api")
    @Tag("negative")
    @DisplayName("GET /logs/{sessionId} without WebSocket upgrade returns 400")
    void sessionLogsRequireWebSocketUpgrade() {
        step("GET /logs/unknown-session without WebSocket headers", () ->
                given()
                        .when()
                        .get("/logs/unknown-session")
                        .then()
                        .statusCode(400)
                        .body(containsString("websocket")));
    }
}
