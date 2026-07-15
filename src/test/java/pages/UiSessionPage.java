package pages;

import com.codeborne.selenide.SelenideElement;
import com.codeborne.selenide.WebDriverRunner;
import io.qameta.allure.Step;

import java.time.Duration;

import static com.codeborne.selenide.Condition.hidden;
import static com.codeborne.selenide.Condition.visible;
import static com.codeborne.selenide.Selenide.$;
import static com.codeborne.selenide.Selenide.webdriver;
import static com.codeborne.selenide.WebDriverConditions.urlContaining;

public class UiSessionPage {

    private static final Duration VNC_CONNECT_TIMEOUT = Duration.ofSeconds(90);

    private final SelenideElement vncCard = $(".vnc-card");
    private final SelenideElement vncConnectionStatus = $(".vnc-connection-status");
    private final SelenideElement lockControl = $(".control_lock");

    @Step("Wait for session page URL")
    public UiSessionPage waitForSessionPage(Duration timeout) {
        webdriver().shouldHave(urlContaining("/sessions/"), timeout);
        return this;
    }

    @Step("Wait until VNC is connected")
    public UiSessionPage shouldVncBeConnected() {
        vncCard.shouldBe(visible, VNC_CONNECT_TIMEOUT);
        vncConnectionStatus.shouldBe(hidden, VNC_CONNECT_TIMEOUT);
        lockControl.shouldBe(visible, VNC_CONNECT_TIMEOUT);
        return this;
    }

    @Step("Unlock VNC screen")
    public UiSessionPage unlockVncScreen() {
        lockControl.click();
        lockControl.$(".dripicons-lock-open").shouldBe(visible);
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
