package tests;

import annotations.Component;
import annotations.Layer;
import api.hub.HubSessionApi;
import com.codeborne.selenide.Condition;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import io.qameta.allure.Story;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.parallel.ResourceAccessMode;
import org.junit.jupiter.api.parallel.ResourceLock;

import java.time.Duration;

import static com.codeborne.selenide.Selenide.$;
import static io.qameta.allure.Allure.step;

@Layer("e2e")
@Component("selenoid-ui")
@Epic("selenoid-ui")
@Feature("UI sessions list")
@Story("UI sessions list")
@DisplayName("UI sessions list")
@ResourceLock(value = "hubSessions", mode = ResourceAccessMode.READ_WRITE)
class UiSessionsListTests extends UiTestBase {

    private static final Duration SESSION_APPEAR_TIMEOUT = Duration.ofSeconds(20);

    @Test
    @Tag("smoke")
    @Tag("positive")
    @DisplayName("Active hub session appears in dashboard sessions list")
    void activeSessionAppearsInSessionsList() {
        var sessionId = step("Create hub session via API", () -> HubSessionApi.create());
        try {
            step("Open dashboard and wait for stack CONNECTED", () ->
                    uiDashboard.openPage().shouldBeConnected());

            step("Verify session row shows browser name", () ->
                    $(".sessions__list .session .browser .name")
                            .shouldHave(Condition.text(config.browser()), SESSION_APPEAR_TIMEOUT));

            step("Verify empty-state message is hidden", () ->
                    $(".no-any").shouldBe(Condition.hidden, SESSION_APPEAR_TIMEOUT));
        } finally {
            step("Delete hub session", () -> HubSessionApi.delete(sessionId));
        }
    }
}
