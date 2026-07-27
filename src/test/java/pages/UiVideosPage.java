package pages;

import com.codeborne.selenide.SelenideElement;
import io.qameta.allure.Step;

import static com.codeborne.selenide.Condition.exist;
import static com.codeborne.selenide.Condition.visible;
import static com.codeborne.selenide.Selenide.$;
import static com.codeborne.selenide.Selenide.open;

/**
 * Finished sessions archive on the Sessions route (replaces the old Videos tab).
 * Route: {@code #/sessions} — Live sessions + Finished sessions panels.
 */
public class UiVideosPage {

    private final SelenideElement archivePanel = $("[data-testid='archive-panel']");
    private final SelenideElement list = $(".archive__list");
    private final SelenideElement pager = $("[data-testid='archive-pager']");
    private final SelenideElement empty = $(".archive-panel .no-any");

    @Step("Open Sessions page (finished sessions archive)")
    public UiVideosPage openPage() {
        open("/#/sessions");
        archivePanel.shouldBe(visible);
        return this;
    }

    @Step("Finished sessions list shell is present")
    public UiVideosPage shouldShowListContainer() {
        // Empty archive keeps `.archive__list` in DOM with zero height (not "visible").
        archivePanel.shouldBe(visible);
        list.should(exist);
        return this;
    }

    @Step("Empty state or pager is present")
    public UiVideosPage shouldShowEmptyOrPager() {
        if (pager.exists()) {
            pager.shouldBe(visible);
        } else {
            empty.shouldBe(visible);
        }
        return this;
    }
}
