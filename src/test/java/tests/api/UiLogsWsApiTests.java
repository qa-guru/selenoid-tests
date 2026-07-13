package tests.api;

import annotations.Component;
import annotations.Layer;
import api.UiApiTestBase;
import api.hub.HubSessionApi;
import api.ui.UiRequest;
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
@Component("selenoid-ui")
@Epic("selenoid-ui")
@Feature("UI logs WebSocket proxy")
@Story("UI logs WebSocket proxy")
@DisplayName("UI /ws/logs WebSocket API")
@ResourceLock(value = "hubSessions", mode = ResourceAccessMode.READ_WRITE)
class UiLogsWsApiTests extends UiApiTestBase {

    @Test
    @Tag("api")
    @Tag("positive")
    @DisplayName("GET /ws/logs/{sessionId} without WebSocket upgrade returns 400 via UI proxy")
    void uiLogsPathRequiresWebSocketUpgrade() {
        var sessionId = step("Create hub session", () -> HubSessionApi.create());
        try {
            step("GET /ws/logs/{sessionId} without WebSocket headers", () ->
                    given()
                            .baseUri(UiRequest.baseUri())
                            .when()
                            .get("/ws/logs/{sessionId}", sessionId)
                            .then()
                            .statusCode(400)
                            .body(containsString("websocket")));
        } finally {
            step("Delete hub session", () -> HubSessionApi.delete(sessionId));
        }
    }
}
