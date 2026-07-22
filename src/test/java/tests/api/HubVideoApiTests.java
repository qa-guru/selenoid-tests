package tests.api;

import annotations.Component;
import annotations.Layer;
import api.ApiTestBase;
import api.hub.HubVideoApi;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import io.qameta.allure.Story;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import static io.qameta.allure.Allure.step;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

@Layer("api")
@Component("video-recorder")
@Epic("video-recorder")
@Feature("Hub video")
@Story("Hub video API")
@DisplayName("Hub video API")
class HubVideoApiTests extends ApiTestBase {

    @Test
    @Tag("api")
    @Tag("positive")
    @DisplayName("GET /video/?json returns paginated page (default limit 10)")
    void videoListJsonReturnsPaginatedPage() {
        var listed = step("GET /video/?json", () -> HubVideoApi.listJson());
        step("Verify paginated payload", () -> {
            assertNotNull(listed);
            assertNotNull(listed.videos());
            assertEquals(10, listed.limit());
            assertEquals(0, listed.offset());
            assertTrue(listed.videos().size() <= listed.limit());
            assertTrue(listed.total() >= listed.videos().size());
        });
    }

    @Test
    @Tag("api")
    @Tag("positive")
    @DisplayName("GET /video/?json&offset=10 returns second page metadata")
    void videoListJsonSecondPage() {
        var first = step("GET first page", () -> HubVideoApi.listJson(10, 0, null));
        var second = step("GET second page", () -> HubVideoApi.listJson(10, 10, null));
        step("Verify second page contract", () -> {
            assertEquals(10, second.limit());
            assertEquals(10, second.offset());
            assertEquals(first.total(), second.total());
            assertTrue(second.videos().size() <= 10);
            if (first.total() <= 10) {
                assertTrue(second.videos().isEmpty());
            }
        });
    }

    @Test
    @Tag("api")
    @Tag("negative")
    @DisplayName("GET /video/{missing}.mp4 returns 404")
    void missingVideoFileReturnsNotFound() {
        step("GET missing video file", () ->
                HubVideoApi.getExpectStatus("missing-session-id.mp4", 404));
    }
}
