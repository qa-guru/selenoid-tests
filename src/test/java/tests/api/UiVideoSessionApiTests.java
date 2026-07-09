package tests.api;

import annotations.Component;
import annotations.Layer;
import api.UiApiTestBase;
import api.hub.HubSessionApi;
import api.ui.UiVideoApi;
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
@Component("selenoid-ui")
@Epic("selenoid-ui")
@Feature("UI session video proxy")
@DisplayName("UI session video via /video proxy")
@ResourceLock(value = "hubSessions", mode = ResourceAccessMode.READ_WRITE)
class UiVideoSessionApiTests extends UiApiTestBase {

    private static final String FULL_CHROME = "148.0";

    @Test
    @Tag("api")
    @Tag("positive")
    @DisplayName("Session video appears in UI GET /video/?json after delete")
    void uiSessionVideoListedAfterClose() throws Exception {
        var sessionId = step("Create hub session with video", () ->
                HubSessionApi.createWithSelenoidOptions(FULL_CHROME, Map.of("enableVideo", true)));
        step("Keep session briefly for recorder", () -> TimeUnit.SECONDS.sleep(3));
        step("Delete hub session", () -> HubSessionApi.delete(sessionId));

        var deadline = System.currentTimeMillis() + TimeUnit.SECONDS.toMillis(30);
        var files = step("Wait for video artifact in UI GET /video/?json", () -> {
            while (System.currentTimeMillis() < deadline) {
                var listed = UiVideoApi.listJson();
                if (listed.stream().anyMatch(name -> name.contains(sessionId) || name.endsWith(".mp4"))) {
                    return listed;
                }
                TimeUnit.SECONDS.sleep(1);
            }
            return UiVideoApi.listJson();
        });

        step("Verify session video is listed via UI proxy", () ->
                assertTrue(files.stream().anyMatch(name -> name.contains(sessionId) || name.endsWith(".mp4")),
                        () -> "expected video for session " + sessionId + " in " + files));
    }
}
