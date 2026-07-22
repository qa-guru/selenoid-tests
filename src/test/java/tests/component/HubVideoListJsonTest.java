package tests.component;

import annotations.Component;
import annotations.Layer;
import api.hub.VideoListResponse;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertTrue;

@Component("video-recorder")
@Layer("component")
@DisplayName("Hub video list fixture")
class HubVideoListJsonTest {

    @Test
    @DisplayName("parses paginated /video/?json page")
    void parsesPaginatedVideoList() {
        var listed = FixtureJson.parse("fixtures/hub/video-list-page.json", VideoListResponse.class);
        assertEquals(10, listed.limit());
        assertEquals(0, listed.offset());
        assertEquals(12, listed.total());
        assertEquals(10, listed.videos().size());
        assertTrue(listed.videos().size() <= listed.limit());
    }
}
