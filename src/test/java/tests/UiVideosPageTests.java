package tests;

import annotations.Component;
import annotations.Layer;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import io.qameta.allure.Story;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;
import pages.UiVideosPage;

import static io.qameta.allure.Allure.step;

@Layer("e2e")
@Component("selenoid-ui")
@Epic("selenoid-ui")
@Feature("UI videos")
@Story("UI videos pagination")
@DisplayName("UI videos page")
class UiVideosPageTests extends UiTestBase {

    private final UiVideosPage videosPage = new UiVideosPage();

    @Test
    @Tag("smoke")
    @Tag("positive")
    @DisplayName("Videos tab loads list container without requiring full catalog")
    void videosTabLoads() {
        step("Open videos page", videosPage::openPage);
        step("Verify list shell", videosPage::shouldShowListContainer);
        step("Verify empty state or pager", videosPage::shouldShowEmptyOrPager);
    }
}
