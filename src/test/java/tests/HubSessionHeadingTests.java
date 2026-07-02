package tests;

import annotations.Layer;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import static com.codeborne.selenide.Condition.text;
import static com.codeborne.selenide.Selenide.$;
import static com.codeborne.selenide.Selenide.open;
import static io.qameta.allure.Allure.step;

@Layer("e2e")
@Epic("selenoid")
@Feature("WebDriver session")
@DisplayName("Hub session heading")
class HubSessionHeadingTests extends TestBase {

    @Test
    @Tag("smoke")
    @Tag("positive")
    @DisplayName("Remote session renders Example Domain heading")
    void remoteSessionRendersHeading() {
        step("Open smoke URL", () -> open(config.smokeUrl()));
        step("Verify heading", () -> $("h1").shouldHave(text("Example Domain")));
    }
}
