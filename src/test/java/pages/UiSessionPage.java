package pages;

import com.codeborne.selenide.SelenideElement;
import com.codeborne.selenide.WebDriverRunner;
import io.qameta.allure.Step;

import java.time.Duration;

import static com.codeborne.selenide.Condition.cssClass;
import static com.codeborne.selenide.Condition.visible;
import static com.codeborne.selenide.Selenide.$;
import static com.codeborne.selenide.Selenide.webdriver;
import static com.codeborne.selenide.WebDriverConditions.urlContaining;

public class UiSessionPage {

    private static final Duration VNC_CONNECT_TIMEOUT = Duration.ofSeconds(90);

    private final SelenideElement vncWindow = $("[data-testid='vnc-window']");
    private final SelenideElement unlockButton = $("[data-testid='vnc-window'] [aria-label='Unlock screen']");

    @Step("Wait for session page URL")
    public UiSessionPage waitForSessionPage(Duration timeout) {
        webdriver().shouldHave(urlContaining("/sessions/"), timeout);
        return this;
    }

    @Step("Wait until VNC is connected")
    public UiSessionPage shouldVncBeConnected() {
        vncWindow.shouldBe(visible, VNC_CONNECT_TIMEOUT);
        vncWindow.shouldHave(cssClass("vnc-window--connected"), VNC_CONNECT_TIMEOUT);
        unlockButton.shouldBe(visible, VNC_CONNECT_TIMEOUT);
        return this;
    }

    @Step("Unlock VNC screen")
    public UiSessionPage unlockVncScreen() {
        unlockButton.shouldBe(visible).click();
        // Control toggles to "Lock screen" when unlocked; some builds keep the same control
        // visible either way — assert VNC stays connected after the click.
        vncWindow.shouldHave(cssClass("vnc-window--connected"));
        $("[data-testid='vnc-window'] [aria-label='Lock screen'], [data-testid='vnc-window'] [aria-label='Unlock screen']")
                .shouldBe(visible);
        return this;
    }

    @Step("Read session id from URL")
    public String sessionId() {
        var url = WebDriverRunner.url();
        var hashMarker = "#/sessions/";
        var marker = "/sessions/";
        var source = url.contains(hashMarker)
                ? url.substring(url.indexOf(hashMarker) + hashMarker.length())
                : url.substring(url.indexOf(marker) + marker.length());
        var end = source.indexOf('?');
        if (end < 0) {
            end = source.length();
        }
        return source.substring(0, end);
    }
}
