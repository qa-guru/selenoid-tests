package tests;

import annotations.Component;
import annotations.Layer;
import com.codeborne.selenide.Configuration;
import com.codeborne.selenide.Selenide;
import com.codeborne.selenide.WebDriverRunner;
import com.microsoft.playwright.Browser;
import com.microsoft.playwright.BrowserContext;
import com.microsoft.playwright.BrowserType;
import com.microsoft.playwright.Page;
import com.microsoft.playwright.Playwright;
import com.microsoft.playwright.options.ServiceWorkerPolicy;
import config.ConfigReader;
import config.TestConfig;
import helpers.HarCapture;
import helpers.HarStats;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import io.qameta.allure.Story;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;
import org.openqa.selenium.MutableCapabilities;
import org.openqa.selenium.chrome.ChromeOptions;
import org.openqa.selenium.json.Json;

import static com.codeborne.selenide.Selenide.closeWebDriver;
import static com.codeborne.selenide.Selenide.open;
import static io.qameta.allure.Allure.step;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * Compares HAR completeness across client capture paths on the same page:
 * <ol>
 *   <li>Native Playwright local {@code recordHar}</li>
 *   <li>Playwright → Selenoid WS {@code recordHar} (same client API, remote browser)</li>
 *   <li>Selenide → Selenoid WebDriver + {@link HarCapture} META (default)</li>
 *   <li>Selenide → Selenoid WebDriver + {@link HarCapture} BODIES (separate session, CDP
 *       {@code getResponseBody}; no hub {@code enableHAR})</li>
 * </ol>
 * Artifacts: {@code build/har-compare/*.har} + {@code summary.json}.
 */
@Layer("e2e")
@Component("selenoid")
@Epic("selenoid")
@Feature("HAR completeness")
@Story("Compare Playwright recordHar vs Selenide HarCapture")
@DisplayName("HAR completeness compare")
@Tag("local-only")
@Tag("har-compare")
class HarCompletenessCompareTests {

    private static final TestConfig config = ConfigReader.testConfig;
    private static final Json JSON = new Json();
    private static final Path OUT = Path.of("build/har-compare");

    /**
     * Best-effort bodies gate on the local fixture: at least one HTTP entry with
     * {@code content.text}. Explicitly <em>not</em> parity with Playwright {@code recordHar}.
     */
    private static final int BODIES_MIN_WITH_CONTENT_TEXT = 1;

    @Test
    @DisplayName("Playwright local / Playwright→Selenoid / Selenide→Selenoid HAR completeness")
    void compareHarCompleteness() throws Exception {
        Files.createDirectories(OUT);

        String localUrl = stripTrailingSlash(ConfigReader.resolveUiUrl());
        String remoteUrl = ConfigReader.resolveUiBrowserUrl();

        List<HarStats> stats = new ArrayList<>();

        step("1) Native Playwright local recordHar → " + localUrl, () ->
                stats.add(capturePlaywrightLocal(localUrl)));

        step("2) Playwright → Selenoid recordHar → " + remoteUrl, () ->
                stats.add(capturePlaywrightSelenoid(remoteUrl)));

        step("3) Selenide → Selenoid HarCapture META (no enableHAR) → " + remoteUrl, () ->
                stats.add(captureSelenideSelenoid(
                        remoteUrl,
                        HarCapture.HarContentMode.META,
                        "3-selenide-HarCapture",
                        "3-selenide-selenoid.har")));

        step("3b) Selenide → Selenoid HarCapture BODIES (no enableHAR) → " + remoteUrl, () ->
                stats.add(captureSelenideSelenoid(
                        remoteUrl,
                        HarCapture.HarContentMode.BODIES,
                        "3b-selenide-HarCapture-bodies",
                        "3b-selenide-HarCapture-bodies.har")));

        HarStats baseline = stats.get(0);
        HarStats pwSelenoid = stats.get(1);
        HarStats selenide = stats.get(2);
        HarStats selenideBodies = stats.get(3);

        String table = HarStats.formatTable(stats);
        System.out.println("\n=== HAR completeness ===\n" + table + "\n");
        System.out.printf(
                "URL coverage vs baseline: pw-selenoid=%.0f%%  selenide(meta)=%.0f%%  "
                        + "selenide(bodies)=%.0f%% withContentText=%d%n",
                pwSelenoid.urlCoverageOf(baseline) * 100,
                selenide.urlCoverageOf(baseline) * 100,
                selenideBodies.urlCoverageOf(baseline) * 100,
                selenideBodies.withContentText);

        Map<String, Object> contentTextNote = new HashMap<>();
        contentTextNote.put(
                "meta",
                "HarCapture default (HarContentMode.META): withContentText==0 on this fixture row");
        contentTextNote.put(
                "bodies",
                "HarCapture HarContentMode.BODIES (live fixture, CDP Network.getResponseBody): "
                        + "withContentText>="
                        + BODIES_MIN_WITH_CONTENT_TEXT
                        + " best-effort; URL cov not weaker than meta; not ≡ recordHar; "
                        + "one writer (no hub enableHAR)");

        Map<String, Object> summary = new HashMap<>();
        summary.put("targetLocal", localUrl);
        summary.put("targetRemote", remoteUrl);
        summary.put("rows", stats.stream().map(HarStats::toRow).toList());
        summary.put("pwSelenoidUrlCoverage", pwSelenoid.urlCoverageOf(baseline));
        summary.put("selenideUrlCoverage", selenide.urlCoverageOf(baseline));
        summary.put("selenideBodiesUrlCoverage", selenideBodies.urlCoverageOf(baseline));
        summary.put("selenideBodiesWithContentText", selenideBodies.withContentText);
        summary.put("bodiesMinWithContentText", BODIES_MIN_WITH_CONTENT_TEXT);
        summary.put("contentTextNote", contentTextNote);
        Files.writeString(OUT.resolve("summary.json"), JSON.toJson(summary), StandardCharsets.UTF_8);
        Files.writeString(OUT.resolve("summary.txt"), table + "\n", StandardCharsets.UTF_8);

        // Meta default contract for the live HarCapture META row.
        assertTrue(selenide.withContentText == 0, "HarCapture meta default must omit content.text");

        // Fair compare: same page, SW blocked — URL set and HTTP entry count within 80% of baseline.
        assertTrue(
                pwSelenoid.urlCoverageOf(baseline) >= 0.80,
                () -> "Playwright→Selenoid URL coverage < 80% of local baseline: "
                        + pwSelenoid.urlCoverageOf(baseline)
                        + "\nbaseline=" + baseline.urls + "\ncandidate=" + pwSelenoid.urls);
        assertTrue(
                selenide.urlCoverageOf(baseline) >= 0.80,
                () -> "Selenide HarCapture URL coverage < 80% of local baseline: "
                        + selenide.urlCoverageOf(baseline)
                        + "\nbaseline=" + baseline.urls + "\ncandidate=" + selenide.urls);
        assertTrue(
                pwSelenoid.httpEntries >= (int) Math.floor(baseline.httpEntries * 0.80),
                () -> "Playwright→Selenoid httpEntries worse than baseline: "
                        + pwSelenoid.httpEntries + " < 80% of " + baseline.httpEntries);
        assertTrue(
                selenide.httpEntries >= (int) Math.floor(baseline.httpEntries * 0.80),
                () -> "Selenide httpEntries worse than baseline: "
                        + selenide.httpEntries + " < 80% of " + baseline.httpEntries);

        // Field richness on shared traffic: status / headers / sizes should not collapse.
        assertTrue(pwSelenoid.withStatus >= pwSelenoid.httpEntries);
        assertTrue(selenide.withStatus >= selenide.httpEntries);
        assertTrue(pwSelenoid.withResponseHeaders >= (int) Math.floor(pwSelenoid.httpEntries * 0.80));
        assertTrue(selenide.withResponseHeaders >= (int) Math.floor(selenide.httpEntries * 0.80));
        assertTrue(pwSelenoid.withContentSize >= (int) Math.floor(pwSelenoid.httpEntries * 0.80));
        assertTrue(selenide.withContentSize >= (int) Math.floor(selenide.httpEntries * 0.80));

        // Bodies opt-in (separate session, one writer): best-effort text gate; URL cov ≥ meta.
        assertTrue(
                selenideBodies.urlCoverageOf(baseline) >= 0.80,
                () -> "HarCapture bodies URL coverage < 80% of local baseline: "
                        + selenideBodies.urlCoverageOf(baseline)
                        + "\nbaseline=" + baseline.urls + "\ncandidate=" + selenideBodies.urls);
        assertTrue(
                selenideBodies.urlCoverageOf(baseline) + 1e-9 >= selenide.urlCoverageOf(baseline),
                () -> "HarCapture bodies URL cov weaker than meta: bodies="
                        + selenideBodies.urlCoverageOf(baseline)
                        + " meta=" + selenide.urlCoverageOf(baseline));
        assertTrue(
                selenideBodies.withContentText >= BODIES_MIN_WITH_CONTENT_TEXT,
                () -> "HarCapture BODIES expected withContentText>="
                        + BODIES_MIN_WITH_CONTENT_TEXT
                        + " on fixture (best-effort), got " + selenideBodies.withContentText);
    }

    private static HarStats capturePlaywrightLocal(String url) throws Exception {
        Path harPath = OUT.resolve("1-playwright-local.har");
        try (Playwright pw = Playwright.create()) {
            Browser browser = pw.chromium().launch(new BrowserType.LaunchOptions().setHeadless(true));
            Browser.NewContextOptions opts = new Browser.NewContextOptions()
                    .setRecordHarPath(harPath)
                    .setServiceWorkers(ServiceWorkerPolicy.BLOCK);
            BrowserContext ctx = browser.newContext(opts);
            Page page = ctx.newPage();
            page.navigate(url);
            page.waitForLoadState();
            page.waitForTimeout(4000);
            ctx.close(); // flush HAR
            browser.close();
        }
        return HarStats.fromFile("1-pw-local-recordHar", harPath);
    }

    private static HarStats capturePlaywrightSelenoid(String url) throws Exception {
        Path harPath = OUT.resolve("2-playwright-selenoid.har");
        try (Playwright pw = Playwright.create()) {
            Browser browser = pw.chromium().connect(ConfigReader.resolvePlaywrightWsEndpoint());
            Browser.NewContextOptions opts = new Browser.NewContextOptions()
                    .setRecordHarPath(harPath)
                    .setServiceWorkers(ServiceWorkerPolicy.BLOCK);
            BrowserContext ctx = browser.newContext(opts);
            Page page = ctx.newPage();
            page.navigate(url);
            page.waitForLoadState();
            page.waitForTimeout(4000);
            try {
                ctx.close();
            } catch (RuntimeException ignored) {
                // hub may already have torn down
            }
            try {
                browser.close();
            } catch (RuntimeException ignored) {
                // ignore
            }
        }
        return HarStats.fromFile("2-pw-selenoid-recordHar", harPath);
    }

    /**
     * Client {@link HarCapture} only — no hub {@code enableHAR} (one writer per session).
     *
     * @param mode     META (default) or BODIES (CDP {@code Network.getResponseBody})
     * @param label    HarStats / MATRIX label
     * @param fileName artifact under {@code build/har-compare/}
     */
    private static HarStats captureSelenideSelenoid(
            String url,
            HarCapture.HarContentMode mode,
            String label,
            String fileName) throws Exception {
        Path harPath = OUT.resolve(fileName);
        Configuration.browser = "chrome";
        Configuration.browserVersion = config.chromeVersion();
        Configuration.browserSize = config.browserSize();
        Configuration.headless = true;
        Configuration.remote = config.remoteUrl();
        Configuration.timeout = 15_000;

        ChromeOptions chrome = new ChromeOptions();
        chrome.addArguments(
                "--headless=new",
                "--no-sandbox",
                "--disable-dev-shm-usage",
                // Match Playwright ServiceWorkerPolicy.BLOCK for fair URL-set compare
                "--disable-features=ServiceWorker");
        HarCapture.enablePerformanceLogging(chrome);

        MutableCapabilities caps = new MutableCapabilities();
        caps.setCapability(ChromeOptions.CAPABILITY, chrome);
        Map<String, Object> selenoid = new HashMap<>();
        selenoid.put("enableVNC", false);
        selenoid.put("enableVideo", false);
        selenoid.put("enableHAR", false);
        selenoid.put("headless", true);
        selenoid.put("name", label);
        caps.setCapability("selenoid:options", selenoid);
        // goog:loggingPrefs must be top-level for remote Chromium
        HarCapture.enablePerformanceLogging(caps);
        Configuration.browserCapabilities = caps;

        try {
            open(url);
            Selenide.sleep(4000);
            byte[] har = HarCapture.collectHarJson(mode)
                    .orElseThrow(() -> new IllegalStateException(
                            "HarCapture " + mode + " produced empty HAR"));
            Files.write(harPath, har);
            return HarStats.fromBytes(label, harPath, har);
        } finally {
            if (WebDriverRunner.hasWebDriverStarted()) {
                closeWebDriver();
            }
        }
    }

    private static String stripTrailingSlash(String url) {
        if (url == null || url.isEmpty()) {
            return url;
        }
        return url.endsWith("/") ? url.substring(0, url.length() - 1) : url;
    }
}
