package tests.api;

import annotations.Component;
import annotations.Layer;
import api.ApiTestBase;
import api.hub.HubRequest;
import api.hub.HubSessionApi;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import io.qameta.allure.Story;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.parallel.ResourceAccessMode;
import org.junit.jupiter.api.parallel.ResourceLock;

import static io.qameta.allure.Allure.step;
import static io.restassured.RestAssured.given;
import static org.hamcrest.Matchers.containsString;

@Layer("api")
@Component("selenoid")
@Epic("selenoid")
@Feature("Hub logs WebSocket")
@Story("Hub logs WebSocket")
@DisplayName("Hub session logs WebSocket API")
@ResourceLock(value = "hubSessions", mode = ResourceAccessMode.READ_WRITE)
class HubLogsSessionApiTests extends ApiTestBase {

    @Test
    @Tag("api")
    @Tag("positive")
    @DisplayName("GET /logs/{sessionId} without WebSocket upgrade returns 400 for active session")
    void sessionLogsRequireWebSocketUpgrade() {
        var sessionId = step("Create hub session", () -> HubSessionApi.create());
        try {
            step("GET /logs/{sessionId} without WebSocket headers", () ->
                    given()
                            .baseUri(HubRequest.baseUri())
                            .when()
                            .get("/logs/{sessionId}", sessionId)
                            .then()
                            .statusCode(400)
                            .body(containsString("websocket")));
        } finally {
            step("Delete hub session", () -> HubSessionApi.delete(sessionId));
        }
    }
}
