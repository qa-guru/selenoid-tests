package tests;

import allure.Attachments;
import com.microsoft.playwright.Browser;
import com.microsoft.playwright.BrowserContext;
import com.microsoft.playwright.Page;
import com.microsoft.playwright.Playwright;
import config.ConfigReader;
import config.TestConfig;
import java.nio.file.Files;
import java.nio.file.Path;
import org.junit.jupiter.api.AfterAll;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.BeforeEach;

public abstract class PlaywrightTestBase {

    protected static final TestConfig config = ConfigReader.testConfig;

    protected static Playwright playwright;
    protected Browser browser;
    protected BrowserContext context;
    protected Page page;
    private Path harPath;

    @BeforeAll
    static void startPlaywright() {
        playwright = Playwright.create();
    }

    @BeforeEach
    void connectBrowser() throws Exception {
        browser = playwright.chromium().connect(ConfigReader.resolvePlaywrightWsEndpoint());
        Browser.NewContextOptions opts = new Browser.NewContextOptions();
        harPath = null;
        if (recordHarEnabled()) {
            harPath = Files.createTempFile("playwright-", ".har");
            opts.setRecordHarPath(harPath);
        }
        context = browser.newContext(opts);
        page = context.newPage();
    }

    @AfterEach
    void closeBrowser() {
        if (context != null) {
            try {
                context.close(); // flushes recordHar to disk
            } catch (RuntimeException ignored) {
                // session may already be gone on hub
            }
            context = null;
        }
        if (recordHarEnabled() && allureResultsEnabled()) {
            Attachments.har(harPath);
        }
        if (harPath != null) {
            try {
                Files.deleteIfExists(harPath);
            } catch (Exception ignored) {
                // temp cleanup best-effort
            }
            harPath = null;
        }
        if (browser != null) {
            try {
                browser.close();
            } catch (RuntimeException ignored) {
                // ignore
            }
            browser = null;
        }
        page = null;
    }

    @AfterAll
    static void stopPlaywright() {
        if (playwright != null) {
            playwright.close();
            playwright = null;
        }
    }

    private static boolean allureResultsEnabled() {
        return !"none".equals(config.allureReportMode());
    }

    private static boolean recordHarEnabled() {
        return config.attachHarLogs();
    }
}
