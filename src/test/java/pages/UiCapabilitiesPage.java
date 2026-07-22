package pages;

import com.codeborne.selenide.SelenideElement;
import io.qameta.allure.Step;

import java.time.Duration;

import static com.codeborne.selenide.Condition.cssClass;
import static com.codeborne.selenide.Condition.enabled;
import static com.codeborne.selenide.Condition.text;
import static com.codeborne.selenide.Condition.visible;
import static com.codeborne.selenide.Selenide.$;
import static com.codeborne.selenide.Selenide.open;

/**
 * Capabilities page for Selenoid 3 UI.
 * Browser choice is {@code PlaqueTagstrip} chips (not react-select).
 */
public class UiCapabilitiesPage {

    private static final Duration BROWSER_LIST_TIMEOUT = Duration.ofSeconds(20);
    private static final Duration CREATE_SESSION_TIMEOUT = Duration.ofSeconds(300);

    private final SelenideElement setupPanel = $("[data-testid=capabilities-setup]");
    private final SelenideElement browserSelect = $("[data-testid=capabilities-browser-select]");
    private final SelenideElement createSessionButton = $("[data-testid=capabilities-create-session]");

    @Step("Open Capabilities page")
    public UiCapabilitiesPage openPage() {
        open("/#/capabilities");
        setupPanel.shouldBe(visible, BROWSER_LIST_TIMEOUT);
        browserSelect.shouldBe(visible, BROWSER_LIST_TIMEOUT);
        browserSelect.$(".plaque-field-seg").shouldBe(visible, BROWSER_LIST_TIMEOUT);
        createSessionButton.shouldBe(visible, BROWSER_LIST_TIMEOUT);
        return this;
    }

    @Step("Wait until Webdriver browser chips are ready")
    public UiCapabilitiesPage openBrowserMenu() {
        // Tagstrip chips are always visible — no react-select menu to expand.
        browserSelect.shouldBe(visible, BROWSER_LIST_TIMEOUT);
        browserSelect.$("button.plaque-field-seg__btn").shouldBe(visible, BROWSER_LIST_TIMEOUT);
        return this;
    }

    @Step("Select chrome {version}")
    public UiCapabilitiesPage selectChrome(String version) {
        openBrowserMenu();
        var chipValue = "chrome_" + version;
        var chip = browserSelect.$("button[data-value='" + chipValue + "']");
        chip.shouldBe(visible, BROWSER_LIST_TIMEOUT).click();
        chip.shouldHave(cssClass("plaque-field-seg__btn--on"));
        createSessionButton.shouldBe(enabled);
        createSessionButton.shouldHave(text("Create Session"));
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
