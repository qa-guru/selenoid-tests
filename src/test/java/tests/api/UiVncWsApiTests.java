package tests.api;

import annotations.Component;
import annotations.Layer;
import api.UiApiTestBase;
import api.hub.HubSessionApi;
import api.ui.UiRequest;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.parallel.ResourceAccessMode;
import org.junit.jupiter.api.parallel.ResourceLock;

import java.util.Map;

import static io.qameta.allure.Allure.step;
import static io.restassured.RestAssured.given;
import static org.hamcrest.Matchers.containsString;

@Layer("api")
@Component("selenoid-ui")
@Epic("selenoid-ui")
@Feature("UI VNC WebSocket proxy")
@DisplayName("UI /ws/vnc WebSocket API")
@ResourceLock(value = "hubSessions", mode = ResourceAccessMode.READ_WRITE)
class UiVncWsApiTests extends UiApiTestBase {

    private static final String FULL_CHROME = "148.0";

    @Test
    @Tag("api")
    @Tag("positive")
    @DisplayName("GET /ws/vnc/{sessionId} without WebSocket upgrade returns 400 via UI proxy")
    void uiVncPathRequiresWebSocketUpgrade() {
        var sessionId = step("Create hub session with VNC", () ->
                HubSessionApi.createWithSelenoidOptions(FULL_CHROME, Map.of("enableVNC", true)));
        try {
            step("GET /ws/vnc/{sessionId} without WebSocket headers", () ->
                    given()
                            .baseUri(UiRequest.baseUri())
                            .when()
                            .get("/ws/vnc/{sessionId}", sessionId)
                            .then()
                            .statusCode(400)
                            .body(containsString("websocket")));
        } finally {
            step("Delete hub session", () -> HubSessionApi.delete(sessionId));
        }
    }
}
