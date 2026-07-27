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

    private final SelenideElement sseStatus = $("#sse-status");
    private final SelenideElement selenoidStatus = $("#selenoid-status");

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

    @Step("Wait until SSE is Connected")
    public UiDashboardPage shouldSseBeConnected() {
        // StatusTile state label is title-case ("Connected"), not ALL CAPS.
        sseStatus.shouldHave(text("Connected"), RECOVERY_TIMEOUT);
        return this;
    }

    @Step("Wait until Selenoid is Connected")
    public UiDashboardPage shouldSelenoidBeConnected() {
        selenoidStatus.shouldHave(text("Connected"), RECOVERY_TIMEOUT);
        return this;
    }

    @Step("Wait until SSE and Selenoid are Connected")
    public UiDashboardPage shouldBeConnected() {
        shouldSseBeConnected();
        shouldSelenoidBeConnected();
        return this;
    }

    @Step("Wait until Selenoid is degraded")
    public UiDashboardPage shouldBeDegraded() {
        selenoidStatus.shouldHave(matchText("Issue|Unknown"), RECOVERY_TIMEOUT);
        return this;
    }

    @Step("Keep Connected stable for {stableMs} ms")
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
