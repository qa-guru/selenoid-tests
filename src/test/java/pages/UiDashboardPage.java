package pages;

import com.codeborne.selenide.SelenideElement;
import io.qameta.allure.Step;

import java.time.Duration;

import static com.codeborne.selenide.Condition.matchText;
import static com.codeborne.selenide.Condition.text;
import static com.codeborne.selenide.Selenide.$;
import static com.codeborne.selenide.Selenide.open;
import static com.codeborne.selenide.Selenide.refresh;

public class UiDashboardPage {

    private static final Duration RECOVERY_TIMEOUT = Duration.ofSeconds(20);

    private final SelenideElement sseStatus = $("[data-testid='sse-status-label']");
    private final SelenideElement selenoidStatus = $("[data-testid='selenoid-status-label']");

    @Step("Open Selenoid UI dashboard")
    public UiDashboardPage openPage() {
        open("/");
        return this;
    }

    @Step("Reload dashboard")
    public UiDashboardPage reloadPage() {
        refresh();
        return this;
    }

    @Step("Wait until SSE is CONNECTED")
    public UiDashboardPage shouldSseBeConnected() {
        sseStatus.shouldHave(text("CONNECTED"), RECOVERY_TIMEOUT);
        return this;
    }

    @Step("Wait until SELENOID is CONNECTED")
    public UiDashboardPage shouldSelenoidBeConnected() {
        selenoidStatus.shouldHave(text("CONNECTED"), RECOVERY_TIMEOUT);
        return this;
    }

    @Step("Wait until SSE and SELENOID are CONNECTED")
    public UiDashboardPage shouldBeConnected() {
        shouldSseBeConnected();
        shouldSelenoidBeConnected();
        return this;
    }

    @Step("Wait until SELENOID is degraded")
    public UiDashboardPage shouldBeDegraded() {
        selenoidStatus.shouldHave(matchText("ISSUE|UNKNOWN"), RECOVERY_TIMEOUT);
        return this;
    }

    @Step("Keep CONNECTED stable for {stableMs} ms")
    public UiDashboardPage shouldStayConnected(long stableMs) throws InterruptedException {
        var stepMs = 500L;
        var steps = stableMs / stepMs;
        for (var i = 0; i < steps; i++) {
            shouldBeConnected();
            Thread.sleep(stepMs);
        }
        return this;
    }
}
