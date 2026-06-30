package tests;

import allure.Attachments;
import com.codeborne.selenide.Configuration;
import com.codeborne.selenide.WebDriverRunner;
import com.codeborne.selenide.logevents.SimpleReport;
import config.ConfigReader;
import config.TestConfig;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.TestInfo;

import static com.codeborne.selenide.Selenide.closeWebDriver;

public class UiTestBase {

    protected static final TestConfig config = ConfigReader.testConfig;
    protected final pages.UiDashboardPage uiDashboard = new pages.UiDashboardPage();
    private static final SimpleReport selenideReport = new SimpleReport();

    @BeforeAll
    static void uiSetup() {
        if (config.logToConsole()) {
            System.setProperty("org.slf4j.simpleLogger.defaultLogLevel", config.rootLogLevel());
        } else {
            System.setProperty("org.slf4j.simpleLogger.defaultLogLevel", "off");
        }

        Configuration.baseUrl = ConfigReader.resolveUiUrl();
        Configuration.browser = config.browser();
        Configuration.browserSize = config.browserSize();
        Configuration.headless = true;
        Configuration.remote = null;
        Configuration.timeout = 20_000;
    }

    @BeforeEach
    void beforeEach() {
        if (config.logToConsole() && config.selenideLogToConsole()) {
            selenideReport.start();
        }
    }

    @AfterEach
    void afterEach(TestInfo testInfo) {
        if (config.logToConsole() && config.selenideLogToConsole()) {
            selenideReport.finish(testInfo.getDisplayName());
        }

        if (config.attachLastScreenshot() && WebDriverRunner.hasWebDriverStarted()) {
            Attachments.screenshot("Last screenshot");
        }

        if (config.closeBrowserAfterEach()) {
            closeWebDriver();
        }
    }
}
