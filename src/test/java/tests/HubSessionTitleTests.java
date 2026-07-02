package tests;

import annotations.Layer;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import static com.codeborne.selenide.Selenide.open;
import static com.codeborne.selenide.WebDriverConditions.title;
import static com.codeborne.selenide.Selenide.webdriver;
import static io.qameta.allure.Allure.step;

@Layer("e2e")
@Epic("selenoid")
@Feature("WebDriver session")
@DisplayName("Hub session title")
class HubSessionTitleTests extends TestBase {

    @Test
    @Tag("smoke")
    @Tag("positive")
    @DisplayName("Remote session page has Example Domain title")
    void remoteSessionHasExampleDomainTitle() {
        step("Open smoke URL", () -> open(config.smokeUrl()));
        step("Verify document title", () -> webdriver().shouldHave(title("Example Domain")));
    }
}
