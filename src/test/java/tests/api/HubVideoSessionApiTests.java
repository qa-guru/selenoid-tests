package tests.api;

import annotations.Component;
import annotations.Layer;
import api.ApiTestBase;
import api.hub.HubSessionApi;
import api.hub.HubVideoApi;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.parallel.ResourceAccessMode;
import org.junit.jupiter.api.parallel.ResourceLock;

import java.util.Map;
import java.util.concurrent.TimeUnit;

import static io.qameta.allure.Allure.step;
import static org.junit.jupiter.api.Assertions.assertTrue;

@Layer("api")
@Component("selenoid")
@Epic("selenoid")
@Feature("Hub session video")
@DisplayName("Hub session video API")
@ResourceLock(value = "hubSessions", mode = ResourceAccessMode.READ_WRITE)
class HubVideoSessionApiTests extends ApiTestBase {

    private static final String FULL_CHROME = "148.0";

    @Test
    @Tag("api")
    @Tag("positive")
    @DisplayName("Session with enableVideo=true appears in GET /video/?json after delete")
    void sessionVideoListedAfterClose() throws Exception {
        var sessionId = step("Create hub session with video", () ->
                HubSessionApi.createWithSelenoidOptions(FULL_CHROME, Map.of("enableVideo", true)));
        step("Keep session briefly for recorder", () -> TimeUnit.SECONDS.sleep(3));
        step("Delete hub session", () -> HubSessionApi.delete(sessionId));

        var deadline = System.currentTimeMillis() + TimeUnit.SECONDS.toMillis(30);
        var files = step("Wait for video artifact in GET /video/?json", () -> {
            while (System.currentTimeMillis() < deadline) {
                var listed = HubVideoApi.listJson();
                if (listed.stream().anyMatch(name -> name.contains(sessionId) || name.endsWith(".mp4"))) {
                    return listed;
                }
                TimeUnit.SECONDS.sleep(1);
            }
            return HubVideoApi.listJson();
        });

        step("Verify session video is listed", () ->
                assertTrue(files.stream().anyMatch(name -> name.contains(sessionId) || name.endsWith(".mp4")),
                        () -> "expected video for session " + sessionId + " in " + files));
    }
}
