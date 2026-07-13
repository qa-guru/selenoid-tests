package tests.api;

import annotations.Component;
import annotations.Layer;
import api.UiApiTestBase;
import api.ui.UiVideoApi;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import io.qameta.allure.Story;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import static io.qameta.allure.Allure.step;
import static org.junit.jupiter.api.Assertions.assertNotNull;

@Layer("api")
@Component("video-recorder")
@Epic("video-recorder")
@Feature("UI video proxy")
@Story("UI /video proxy API")
@DisplayName("UI /video proxy API")
class UiVideoApiTests extends UiApiTestBase {

    @Test
    @Tag("api")
    @Tag("positive")
    @DisplayName("GET /video/?json via UI proxy returns file name list")
    void uiVideoListJsonReturnsArray() {
        var files = step("GET UI /video/?json", UiVideoApi::listJson);
        step("Verify list payload", () -> assertNotNull(files));
    }

    @Test
    @Tag("api")
    @Tag("negative")
    @DisplayName("GET /video/{missing}.mp4 via UI proxy returns 404")
    void uiMissingVideoFileReturnsNotFound() {
        step("GET missing video via UI", () ->
                UiVideoApi.getExpectStatus("missing-session-id.mp4", 404));
    }
}
