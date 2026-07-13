package tests;

import annotations.Component;
import annotations.Layer;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import io.qameta.allure.Story;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import static com.codeborne.selenide.Condition.text;
import static com.codeborne.selenide.Selenide.$;
import static com.codeborne.selenide.Selenide.open;
import static com.codeborne.selenide.Selenide.sessionId;
import static com.codeborne.selenide.Selenide.webdriver;
import static com.codeborne.selenide.WebDriverConditions.title;
import static io.qameta.allure.Allure.step;
import static org.junit.jupiter.api.Assertions.assertFalse;

@Layer("e2e")
@Component("webdriver-image")
@Epic("webdriver-image")
@Feature("WebDriver session")
@Story("WebDriver session")
@DisplayName("WebDriver session")
class HubSessionTests extends TestBase {

    @Test
    @Tag("smoke")
    @Tag("positive")
    @DisplayName("Remote Chrome session opens example.com")
    void remoteSessionOpensExampleDomain() {
        step("Open smoke URL via remote WebDriver", () -> open(config.smokeUrl()));

        step("Verify session id is assigned", () ->
                assertFalse(sessionId().toString().isBlank()));

        step("Verify page title", () -> webdriver().shouldHave(title("Example Domain")));

        step("Verify heading text", () -> $("h1").shouldHave(text("Example Domain")));
    }
}
