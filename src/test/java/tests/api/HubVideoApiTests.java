package tests.api;

import annotations.Component;
import annotations.Layer;
import api.ApiTestBase;
import api.hub.HubVideoApi;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import static io.qameta.allure.Allure.step;
import static org.junit.jupiter.api.Assertions.assertNotNull;

@Layer("api")
@Component("video-recorder")
@Epic("video-recorder")
@Feature("Hub video")
@DisplayName("Hub video API")
class HubVideoApiTests extends ApiTestBase {

    @Test
    @Tag("api")
    @Tag("positive")
    @DisplayName("GET /video/?json returns file name list")
    void videoListJsonReturnsArray() {
        var files = step("GET /video/?json", HubVideoApi::listJson);
        step("Verify list payload", () -> assertNotNull(files));
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
