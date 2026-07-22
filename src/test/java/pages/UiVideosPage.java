package pages;

import com.codeborne.selenide.SelenideElement;
import io.qameta.allure.Step;

import static com.codeborne.selenide.Condition.visible;
import static com.codeborne.selenide.Selenide.$;
import static com.codeborne.selenide.Selenide.open;

public class UiVideosPage {

    private final SelenideElement list = $(".videos__list");
    private final SelenideElement pager = $("[data-testid='videos-pager']");
    private final SelenideElement empty = $(".no-any");

    @Step("Open Selenoid UI videos page")
    public UiVideosPage openPage() {
        open("/#/videos");
        return this;
    }

    @Step("Videos list container is visible")
    public UiVideosPage shouldShowListContainer() {
        list.shouldBe(visible);
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
