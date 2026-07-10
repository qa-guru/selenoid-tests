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

import java.util.List;
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

    @Test
    @Tag("api")
    @Tag("positive")
    @DisplayName("Session video is downloadable via UI /video proxy after delete")
    void uiSessionVideoListedAfterClose() throws Exception {
        var sessionId = step("Create hub session with video", () ->
                HubSessionApi.createWithSelenoidOptions(config.browserVersion(), Map.of("enableVideo", true)));
        step("Keep session briefly for recorder", () -> TimeUnit.SECONDS.sleep(3));
        step("Delete hub session", () -> HubSessionApi.delete(sessionId));

        var deadline = System.currentTimeMillis() + TimeUnit.SECONDS.toMillis(30);
        var videoFile = step("Wait for video artifact in UI GET /video/?json", () -> {
            while (System.currentTimeMillis() < deadline) {
                var listed = UiVideoApi.listJson();
                var match = findSessionVideo(listed, sessionId);
                if (match != null) {
                    return match;
                }
                TimeUnit.SECONDS.sleep(1);
            }
            return findSessionVideo(UiVideoApi.listJson(), sessionId);
        });

        step("Verify session video is listed via UI proxy", () ->
                assertTrue(videoFile != null,
                        () -> "expected video for session " + sessionId + " in UI video list"));

        var body = step("Download session video via UI proxy", () -> UiVideoApi.download(videoFile));
        step("Verify MP4 payload via UI proxy", () -> assertValidMp4(body, videoFile));
    }

    private static String findSessionVideo(List<String> files, String sessionId) {
        return files.stream()
                .filter(name -> name.contains(sessionId))
                .findFirst()
                .orElse(null);
    }

    private static void assertValidMp4(byte[] body, String fileName) {
        assertTrue(body.length > 1024,
                () -> "expected non-trivial mp4 body for " + fileName + ", got " + body.length + " bytes");
        assertTrue(body.length >= 8
                        && body[4] == 'f'
                        && body[5] == 't'
                        && body[6] == 'y'
                        && body[7] == 'p',
                () -> "expected ISO BMFF ftyp box in " + fileName);
    }
}
