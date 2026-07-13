package tests.api;

import annotations.Component;
import annotations.Layer;
import api.UiApiTestBase;
import api.ui.UiRequest;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import io.qameta.allure.Story;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import static io.qameta.allure.Allure.step;
import static io.restassured.RestAssured.given;
import static org.hamcrest.Matchers.containsString;

@Layer("api")
@Component("selenoid-ui")
@Epic("selenoid-ui")
@Feature("UI VNC WebSocket proxy")
@Story("UI /ws/vnc WebSocket API")
@DisplayName("UI /ws/vnc WebSocket API")
class UiVncWsApiTests extends UiApiTestBase {

    @Test
    @Tag("api")
    @Tag("positive")
    @DisplayName("GET /ws/vnc/{sessionId} without WebSocket upgrade returns 400 via UI proxy")
    void uiVncPathRequiresWebSocketUpgrade() {
        step("GET /ws/vnc/{sessionId} without WebSocket headers", () ->
                given()
                        .baseUri(UiRequest.baseUri())
                        .when()
                        .get("/ws/vnc/{sessionId}", "abc")
                        .then()
                        .statusCode(400)
                        .body(containsString("websocket")));
    }
}
