package pages;

import com.codeborne.selenide.SelenideElement;
import com.codeborne.selenide.WebDriverRunner;
import config.ConfigReader;
import io.qameta.allure.Step;

import java.time.Duration;

import static com.codeborne.selenide.Condition.cssClass;
import static com.codeborne.selenide.Condition.enabled;
import static com.codeborne.selenide.Condition.text;
import static com.codeborne.selenide.Condition.visible;
import static com.codeborne.selenide.Selenide.$;
import static com.codeborne.selenide.Selenide.open;
import static com.codeborne.selenide.Selenide.sleep;

/**
 * New Session page for Selenoid 3 UI (former Capabilities).
 * Browser choice is {@code PlaqueTagstrip} chips (not react-select).
 */
public class UiCapabilitiesPage {

    private static final Duration BROWSER_LIST_TIMEOUT = Duration.ofSeconds(20);
    private static final Duration CREATE_NAVIGATE_TIMEOUT = Duration.ofSeconds(45);
    private static final Duration CREATE_FALLBACK_TIMEOUT = Duration.ofSeconds(60);

    private final SelenideElement setupPanel = $("[data-testid=capabilities-setup]");
    private final SelenideElement browserSelect = $("[data-testid=capabilities-browser-select]");
    private final SelenideElement createSessionButton = $("[data-testid=capabilities-create-session]");
    private final SelenideElement authUser = $("[data-testid=capabilities-caps-auth-user]");
    private final SelenideElement authPass = $("[data-testid=capabilities-caps-auth-pass]");
    private final SelenideElement enableVideoFalse =
            $("[data-testid=caps-enable-video] button[data-value='false']");

    @Step("Open New Session page")
    public UiCapabilitiesPage openPage() {
        open("/#/new-session");
        setupPanel.shouldBe(visible, BROWSER_LIST_TIMEOUT);
        browserSelect.shouldBe(visible, BROWSER_LIST_TIMEOUT);
        browserSelect.$(".plaque-field-seg").shouldBe(visible, BROWSER_LIST_TIMEOUT);
        createSessionButton.shouldBe(visible, BROWSER_LIST_TIMEOUT);
        return this;
    }

    @Step("Ensure hub Basic Auth fields are filled for Create Session")
    public UiCapabilitiesPage ensureHubAuthFilled() {
        // Remote hub panel (auth fields) mounts only after a browser chip is selected.
        authUser.shouldBe(visible, BROWSER_LIST_TIMEOUT);
        var creds = ConfigReader.resolveHubBasicAuth();
        if (creds[0].isBlank()) {
            return this;
        }
        if (authUser.getValue() == null || authUser.getValue().isBlank()) {
            authUser.setValue(creds[0]);
        }
        if (authPass.getValue() == null || authPass.getValue().isBlank()) {
            authPass.setValue(creds[1]);
        }
        return this;
    }

    @Step("Disable enableVideo for faster smoke create")
    public UiCapabilitiesPage disableVideo() {
        enableVideoFalse.shouldBe(visible, BROWSER_LIST_TIMEOUT).click();
        enableVideoFalse.shouldHave(cssClass("plaque-field-seg__btn--on"));
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
        ensureHubAuthFilled();
        disableVideo();
        createSessionButton.shouldBe(enabled);
        createSessionButton.shouldHave(text("Create Session"));
        return this;
    }

    @Step("Create session from New Session page")
    public UiSessionPage createSession() {
        ensureHubAuthFilled();
        createSessionButton.click();

        var deadline = System.currentTimeMillis() + CREATE_NAVIGATE_TIMEOUT.toMillis();
        while (System.currentTimeMillis() < deadline) {
            if (WebDriverRunner.url().contains("/sessions/")) {
                return new UiSessionPage();
            }
            var title = createSessionButton.getAttribute("title");
            if (title != null && !title.isBlank() && createSessionButton.is(enabled)) {
                throw new AssertionError("Create Session failed: " + title);
            }
            sleep(500);
        }

        // Prod ≤v3.0.11: create can succeed while post-create window/rect hangs without Basic Auth,
        // leaving the UI on #/new-session. Recover via Live sessions list (Manual session).
        open("/#/sessions");
        var sessionRow = $(".sessions__list .session").shouldBe(visible, CREATE_FALLBACK_TIMEOUT);
        var named = $(".sessions__list .session .session-name[title='Manual session']");
        if (named.exists()) {
            sessionRow = named.closest(".session");
        }
        var href = sessionRow.$("a.link.id").shouldBe(visible).getAttribute("href");
        if (href == null || !href.contains("/sessions/")) {
            throw new AssertionError("Live sessions fallback: missing session href, got: " + href);
        }
        var id = href.substring(href.indexOf("/sessions/") + "/sessions/".length());
        var q = id.indexOf('?');
        if (q >= 0) {
            id = id.substring(0, q);
        }
        open("/#/sessions/" + id);
        return new UiSessionPage().waitForSessionPage(CREATE_FALLBACK_TIMEOUT);
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
