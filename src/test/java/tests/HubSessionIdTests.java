package tests;

import annotations.Component;
import annotations.Layer;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import io.qameta.allure.Story;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import static com.codeborne.selenide.Selenide.open;
import static com.codeborne.selenide.Selenide.sessionId;
import static io.qameta.allure.Allure.step;
import static org.junit.jupiter.api.Assertions.assertFalse;

@Layer("e2e")
@Component("webdriver-image")
@Epic("webdriver-image")
@Feature("WebDriver session")
@Story("Hub session id")
@DisplayName("Hub session id")
class HubSessionIdTests extends TestBase {

    @Test
    @Tag("smoke")
    @Tag("positive")
    @DisplayName("Remote session assigns session id")
    void remoteSessionAssignsSessionId() {
        step("Open smoke URL", () -> open(config.smokeUrl()));
        step("Verify session id", () -> assertFalse(sessionId().toString().isBlank()));
    }
}
