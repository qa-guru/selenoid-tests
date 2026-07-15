package pages;

import com.codeborne.selenide.SelenideElement;
import io.qameta.allure.Step;

import java.time.Duration;

import static com.codeborne.selenide.Condition.enabled;
import static com.codeborne.selenide.Condition.text;
import static com.codeborne.selenide.Condition.visible;
import static com.codeborne.selenide.Selenide.$;
import static com.codeborne.selenide.Selenide.open;

public class UiCapabilitiesPage {

    private static final Duration BROWSER_LIST_TIMEOUT = Duration.ofSeconds(20);
    private static final Duration CREATE_SESSION_TIMEOUT = Duration.ofSeconds(300);

    private final SelenideElement setupPanel = $("[data-testid=capabilities-setup]");
    private final SelenideElement browserSelect = $(".capabilities-browser-select");
    private final SelenideElement createSessionButton = $("[data-testid=capabilities-create-session]");

    @Step("Open Capabilities page")
    public UiCapabilitiesPage openPage() {
        open("/#/capabilities");
        setupPanel.shouldBe(visible, BROWSER_LIST_TIMEOUT);
        browserSelect.shouldBe(visible, BROWSER_LIST_TIMEOUT);
        createSessionButton.shouldBe(visible, BROWSER_LIST_TIMEOUT);
        return this;
    }

    @Step("Open browser select menu")
    public UiCapabilitiesPage openBrowserMenu() {
        browserSelect.shouldBe(enabled, BROWSER_LIST_TIMEOUT);
        browserSelect.$(".Select__control").click();
        browserSelect.$(".Select__menu").shouldBe(visible, BROWSER_LIST_TIMEOUT);
        return this;
    }

    @Step("Select chrome {version}")
    public UiCapabilitiesPage selectChrome(String version) {
        openBrowserMenu();
        browserSelect.$$(".Select__option")
                .findBy(text("chrome: " + version))
                .click();
        createSessionButton.shouldBe(enabled);
        return this;
    }

    @Step("Create session from Capabilities")
    public UiSessionPage createSession() {
        createSessionButton.click();
        return new UiSessionPage().waitForSessionPage(CREATE_SESSION_TIMEOUT);
    }

    public SelenideElement setupPanel() {
        return setupPanel;
    }

    public SelenideElement browserSelect() {
        return browserSelect;
    }

    public SelenideElement createSessionButton() {
        return createSessionButton;
    }
}
