package tests.api;

import annotations.Component;
import annotations.Layer;
import api.ApiTestBase;
import api.hub.HubRequest;
import api.hub.HubSessionApi;
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
@Component("selenoid")
@Epic("selenoid")
@Feature("Hub VNC WebSocket")
@DisplayName("Hub session VNC WebSocket API")
@ResourceLock(value = "hubSessions", mode = ResourceAccessMode.READ_WRITE)
class HubVncSessionApiTests extends ApiTestBase {

    @Test
    @Tag("api")
    @Tag("positive")
    @DisplayName("GET /vnc/{sessionId} without WebSocket upgrade returns 400 when enableVNC=true")
    void sessionVncPathRequiresWebSocketUpgrade() {
        var sessionId = step("Create hub session with VNC", () ->
                HubSessionApi.createWithSelenoidOptions(config.browserVersion(), Map.of("enableVNC", true)));
        try {
            step("GET /vnc/{sessionId} without WebSocket headers", () ->
                    given()
                            .baseUri(HubRequest.baseUri())
                            .when()
                            .get("/vnc/{sessionId}", sessionId)
                            .then()
                            .statusCode(400)
                            .body(containsString("websocket")));
        } finally {
            step("Delete hub session", () -> HubSessionApi.delete(sessionId));
        }
    }
}
