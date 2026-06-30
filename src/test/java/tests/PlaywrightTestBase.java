package tests;

import com.microsoft.playwright.Browser;
import com.microsoft.playwright.Page;
import com.microsoft.playwright.Playwright;
import config.ConfigReader;
import config.TestConfig;
import org.junit.jupiter.api.AfterAll;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.BeforeEach;

public abstract class PlaywrightTestBase {

    protected static final TestConfig config = ConfigReader.testConfig;

    protected static Playwright playwright;
    protected Browser browser;
    protected Page page;

    @BeforeAll
    static void startPlaywright() {
        playwright = Playwright.create();
    }

    @BeforeEach
    void connectBrowser() {
        browser = playwright.chromium().connect(ConfigReader.resolvePlaywrightWsEndpoint());
        page = browser.newPage();
    }

    @AfterEach
    void closeBrowser() {
        if (browser != null) {
            browser.close();
            browser = null;
        }
    }

    @AfterAll
    static void stopPlaywright() {
        if (playwright != null) {
            playwright.close();
            playwright = null;
        }
    }
}
