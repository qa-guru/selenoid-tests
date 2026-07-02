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
import org.openqa.selenium.MutableCapabilities;

import java.util.Map;

import static com.codeborne.selenide.Selenide.closeWebDriver;

public class TestBase {

    protected static final TestConfig config = ConfigReader.testConfig;
    private static final SimpleReport selenideReport = new SimpleReport();

    @BeforeAll
    static void setup() {
        if (config.logToConsole()) {
            System.setProperty("org.slf4j.simpleLogger.defaultLogLevel", config.rootLogLevel());
        } else {
            System.setProperty("org.slf4j.simpleLogger.defaultLogLevel", "off");
        }

        Configuration.browser = config.browser();
        Configuration.browserVersion = config.browserVersion();
        Configuration.browserSize = config.browserSize();
        Configuration.headless = config.headless();
        Configuration.remote = config.remoteUrl();

        var capabilities = new MutableCapabilities();
        capabilities.setCapability("goog:chromeOptions", java.util.Map.of(
                "args", java.util.List.of(
                        "--headless=new",
                        "--no-sandbox",
                        "--disable-dev-shm-usage"
                )
        ));
        capabilities.setCapability("selenoid:options", Map.of(
                "enableVNC", config.enableVnc(),
                "enableVideo", config.enableVideo(),
                "headless", config.headless()
        ));
        Configuration.browserCapabilities = capabilities;
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
