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

import java.util.List;
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
    @DisplayName("Session with enableVideo=true produces downloadable MP4 after delete")
    void sessionVideoListedAfterClose() throws Exception {
        var sessionId = step("Create hub session with video", () ->
                HubSessionApi.createWithSelenoidOptions(FULL_CHROME, Map.of("enableVideo", true)));
        step("Keep session briefly for recorder", () -> TimeUnit.SECONDS.sleep(3));
        step("Delete hub session", () -> HubSessionApi.delete(sessionId));

        var deadline = System.currentTimeMillis() + TimeUnit.SECONDS.toMillis(30);
        var videoFile = step("Wait for video artifact in GET /video/?json", () -> {
            while (System.currentTimeMillis() < deadline) {
                var listed = HubVideoApi.listJson();
                var match = findSessionVideo(listed, sessionId);
                if (match != null) {
                    return match;
                }
                TimeUnit.SECONDS.sleep(1);
            }
            return findSessionVideo(HubVideoApi.listJson(), sessionId);
        });

        step("Verify session video is listed", () ->
                assertTrue(videoFile != null,
                        () -> "expected video for session " + sessionId + " in hub video list"));

        var body = step("Download session video", () -> HubVideoApi.download(videoFile));
        step("Verify MP4 payload", () -> assertValidMp4(body, videoFile));
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
