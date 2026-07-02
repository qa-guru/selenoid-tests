package tests;

import annotations.Component;
import annotations.Layer;
import com.codeborne.selenide.Configuration;
import config.ConfigReader;
import helpers.CmInstallerHelper;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import org.junit.jupiter.api.AfterAll;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.Timeout;
import org.junit.jupiter.api.parallel.Execution;
import org.junit.jupiter.api.parallel.ExecutionMode;

import static com.codeborne.selenide.Condition.text;
import static com.codeborne.selenide.Selenide.$;
import static com.codeborne.selenide.Selenide.open;
import static com.codeborne.selenide.Selenide.sessionId;
import static com.codeborne.selenide.Selenide.webdriver;
import static com.codeborne.selenide.WebDriverConditions.title;
import static io.qameta.allure.Allure.step;
import static org.junit.jupiter.api.Assertions.assertFalse;

@Layer("e2e")
@Component("cm")
@Epic("cm")
@Feature("Installed stack session")
@DisplayName("CM installer session")
@Tag("smoke")
@Tag("cm")
@Tag("local-only")
@Execution(ExecutionMode.SAME_THREAD)
@Timeout(300)
class CmInstallerSessionTests extends TestBase {

    private static CmInstallerHelper installer;

    @BeforeAll
    static void installStackViaCm() throws Exception {
        installer = CmInstallerHelper.withTempConfigDir();
        installer.stopAll();

        step("Configure CM-managed stack", () ->
                installer.configure().requireSuccess("configure"));
        step("Start hub via cm", () ->
                installer.startHub().requireSuccess("start hub"));
        step("Wait for hub readiness", () -> installer.waitForHubReady(60_000));
        step("Point Selenide at CM hub", () ->
                Configuration.remote = ConfigReader.resolveCmRemoteUrl());
    }

    @AfterAll
    static void tearDownStack() {
        if (installer != null) {
            installer.stopAll();
            installer.deleteConfigDir();
            installer = null;
        }
    }

    @Test
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
