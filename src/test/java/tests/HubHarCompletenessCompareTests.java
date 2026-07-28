package tests;

import annotations.Component;
import annotations.Layer;
import api.hub.HubHarApi;
import api.hub.HubSessionApi;
import api.hub.HubStatusApi;
import com.microsoft.playwright.Browser;
import com.microsoft.playwright.BrowserContext;
import com.microsoft.playwright.BrowserType;
import com.microsoft.playwright.Page;
import com.microsoft.playwright.Playwright;
import com.microsoft.playwright.options.ServiceWorkerPolicy;
import config.ConfigReader;
import config.TestConfig;
import helpers.HarStats;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import io.qameta.allure.Story;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.time.Duration;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.HashSet;
import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.concurrent.TimeUnit;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;
import org.openqa.selenium.json.Json;

import static io.qameta.allure.Allure.step;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * Hub CDP HAR ({@code enableHAR}) vs client-side baseline from
 * {@link HarCompletenessCompareTests}. Artifacts under {@code build/har-compare/}.
 *
 * <p>One writer per session: hub {@code enableHAR} only — no client {@code recordHar}/HarCapture.
 */
@Layer("e2e")
@Component("selenoid")
@Epic("selenoid")
@Feature("HAR completeness")
@Story("Compare hub enableHAR vs Playwright local recordHar baseline")
@DisplayName("Hub HAR completeness compare")
@Tag("local-only")
@Tag("har-compare")
@Tag("hub-har")
class HubHarCompletenessCompareTests {

    private static final TestConfig config = ConfigReader.testConfig;
    private static final Json JSON = new Json();
    private static final Path OUT = Path.of("build/har-compare");

    @Test
    @DisplayName("Hub WD + PW enableHAR vs 1-pw-local-recordHar baseline")
    void compareHubHarCompleteness() throws Exception {
        Files.createDirectories(OUT);

        String localUrl = stripTrailingSlash(ConfigReader.resolveUiUrl());
        String remoteUrl = ConfigReader.resolveUiBrowserUrl();

        List<HarStats> stats = new ArrayList<>();
        List<String> gaps = new ArrayList<>();

        HarStats baseline = step("0) Baseline 1-pw-local-recordHar", () -> loadOrCaptureBaseline(localUrl));
        stats.add(baseline);

        HarStats wdHub = step("4) WebDriver → Selenoid hub enableHAR → " + remoteUrl, () ->
                captureWebDriverHubHar(remoteUrl));
        stats.add(wdHub);

        HarStats pwHub = step("5) Playwright → Selenoid hub enableHAR (no client recordHar) → " + remoteUrl, () ->
                capturePlaywrightHubHar(remoteUrl, gaps));
        if (pwHub != null) {
            stats.add(pwHub);
        }

        String table = HarStats.formatTable(stats);
        System.out.println("\n=== Hub HAR completeness ===\n" + table + "\n");
        System.out.printf(
                "URL coverage vs baseline: wd-hub=%.0f%%%n",
                wdHub.urlCoverageOf(baseline) * 100);
        if (pwHub != null) {
            System.out.printf("URL coverage vs baseline: pw-hub=%.0f%%%n", pwHub.urlCoverageOf(baseline) * 100);
        }
        if (!gaps.isEmpty()) {
            System.out.println("Gaps:\n- " + String.join("\n- ", gaps));
        }

        // Document known field gaps vs client recordHar.
        gaps.add("hub CDP HAR omits content.text — by design in selenoid/har (size/mimeType only)");
        gaps.add("hub CDP HAR often has content.size=0 (no Network.getResponseBody) and status=0 for in-flight requests at quit");
        gaps.add("do not combine hub enableHAR with client recordHar/HarCapture on one session");

        Map<String, Object> summary = new HashMap<>();
        summary.put("targetLocal", localUrl);
        summary.put("targetRemote", remoteUrl);
        summary.put("rows", stats.stream().map(HarStats::toRow).toList());
        summary.put("wdHubUrlCoverage", wdHub.urlCoverageOf(baseline));
        if (pwHub != null) {
            summary.put("pwHubUrlCoverage", pwHub.urlCoverageOf(baseline));
        } else {
            summary.put("pwHubUrlCoverage", 0.0);
            summary.put("pwHubSkipped", true);
        }
        summary.put("gaps", gaps);
        summary.put("contentTextNote", "hub CDP HAR: withContentText expected ~0; client recordHar fills content.text");
        Files.writeString(OUT.resolve("hub-summary.json"), JSON.toJson(summary), StandardCharsets.UTF_8);
        Files.writeString(
                OUT.resolve("hub-summary.txt"),
                table + "\n\nGaps:\n- " + String.join("\n- ", gaps) + "\n",
                StandardCharsets.UTF_8);

        // Primary gate (product criterion): URL set coverage ≥ 80% of client baseline.
        assertTrue(
                wdHub.urlCoverageOf(baseline) >= 0.80,
                () -> "WebDriver hub-HAR URL coverage < 80% of local baseline: "
                        + wdHub.urlCoverageOf(baseline)
                        + "\nbaseline=" + baseline.urls + "\ncandidate=" + wdHub.urls);
        assertTrue(
                wdHub.httpEntries >= (int) Math.floor(baseline.httpEntries * 0.80),
                () -> "WebDriver hub-HAR httpEntries worse than baseline: "
                        + wdHub.httpEntries + " < 80% of " + baseline.httpEntries);
        // Hub CDP path: content.text always empty; status/size may be 0 for in-flight
        // entries cancelled at session quit — documented in hub-summary gaps, not a hard fail.
        assertTrue(wdHub.withContentText == 0, "hub HAR must not include content.text");
        assertTrue(wdHub.withRequestHeaders >= (int) Math.floor(wdHub.httpEntries * 0.80));

        assertTrue(pwHub != null, "Playwright hub-HAR must be produced (Chromium image :7070)");
        assertTrue(
                pwHub.urlCoverageOf(baseline) >= 0.80,
                () -> "Playwright hub-HAR URL coverage < 80% of local baseline: "
                        + pwHub.urlCoverageOf(baseline)
                        + "\nbaseline=" + baseline.urls + "\ncandidate=" + pwHub.urls);
    }

    private static HarStats loadOrCaptureBaseline(String url) throws Exception {
        Path harPath = OUT.resolve("1-playwright-local.har");
        if (Files.isRegularFile(harPath) && Files.size(harPath) > 100) {
            return HarStats.fromFile("1-pw-local-recordHar", harPath);
        }
        try (Playwright pw = Playwright.create()) {
            Browser browser = pw.chromium().launch(new BrowserType.LaunchOptions().setHeadless(true));
            Browser.NewContextOptions opts = new Browser.NewContextOptions()
                    .setRecordHarPath(harPath)
                    .setServiceWorkers(ServiceWorkerPolicy.BLOCK);
            BrowserContext ctx = browser.newContext(opts);
            Page page = ctx.newPage();
            page.navigate(url);
            page.waitForLoadState();
            page.waitForTimeout(1500);
            ctx.close();
            browser.close();
        }
        return HarStats.fromFile("1-pw-local-recordHar", harPath);
    }

    private static HarStats captureWebDriverHubHar(String url) throws Exception {
        Path harPath = OUT.resolve("4-wd-hub-enableHAR.har");
        Map<String, Object> selenoid = new HashMap<>();
        selenoid.put("enableVNC", true);
        selenoid.put("enableVideo", false);
        selenoid.put("enableHAR", true);
        selenoid.put("enableLog", false);
        selenoid.put("sessionTimeout", "2m");
        selenoid.put("name", "hub-har-wd-compare");

        String sessionId = HubSessionApi.createWithSelenoidOptions(config.chromeVersion(), selenoid);
        try {
            HubSessionApi.navigate(sessionId, url);
            TimeUnit.MILLISECONDS.sleep(3000);
        } finally {
            HubSessionApi.delete(sessionId);
        }

        String file = HubHarApi.waitForSessionHar(sessionId, Duration.ofSeconds(20));
        assertTrue(file != null, () -> "expected hub HAR for WD session " + sessionId);
        byte[] body = HubHarApi.download(file);
        Files.write(harPath, body);
        return HarStats.fromBytes("4-wd-hub-enableHAR", harPath, body);
    }

    /**
     * Hub enableHAR on Playwright WS — no client recordHar.
     * Requires Chromium-family images with DevTools proxy on :7070.
     */
    private static HarStats capturePlaywrightHubHar(String url, List<String> gaps) throws Exception {
        Path harPath = OUT.resolve("5-pw-hub-enableHAR.har");
        String base = ConfigReader.resolvePlaywrightWsEndpoint();
        String sep = base.contains("?") ? "&" : "?";
        String ws = base + sep + "enableHAR=true&enableVideo=false&enableVNC=false&name=hub-har-pw-compare";

        Set<String> before = sessionIds();
        String sessionId;
        try (Playwright pw = Playwright.create()) {
            Browser browser = pw.chromium().connect(ws);
            sessionId = waitNewSessionId(before, Duration.ofSeconds(30));
            BrowserContext ctx = browser.newContext(new Browser.NewContextOptions()
                    .setServiceWorkers(ServiceWorkerPolicy.BLOCK));
            Page page = ctx.newPage();
            // Hub attaches to /page asynchronously after newPage(); yield so
            // Network.enable wins the race before the first navigation.
            page.waitForTimeout(750);
            page.navigate(url);
            page.waitForLoadState();
            page.waitForTimeout(1500);
            try {
                ctx.close();
            } catch (RuntimeException ignored) {
                // hub may tear down with the browser
            }
            try {
                browser.close();
            } catch (RuntimeException ignored) {
                // ignore
            }
        }

        if (sessionId == null) {
            gaps.add("Playwright hub session id not observed in /status — cannot fetch /har");
            return null;
        }

        String file = HubHarApi.waitForSessionHar(sessionId, Duration.ofSeconds(20));
        if (file == null) {
            gaps.add("Playwright hub-HAR missing for " + sessionId
                    + " — expected DevTools :7070 on playwright-chromium/chrome/msedge"
                    + " (HAR_CAPTURE_FAILED)");
            return null;
        }
        byte[] body = HubHarApi.download(file);
        Files.write(harPath, body);
        return HarStats.fromBytes("5-pw-hub-enableHAR", harPath, body);
    }

    private static Set<String> sessionIds() {
        Set<String> ids = new HashSet<>();
        collectSessionIds(HubStatusApi.fetch().browsers(), ids);
        return ids;
    }

    @SuppressWarnings("unchecked")
    private static void collectSessionIds(Object node, Set<String> ids) {
        if (node instanceof Map<?, ?> map) {
            Object id = map.get("id");
            if (id instanceof String s && !s.isBlank()) {
                ids.add(s);
            }
            for (Object v : map.values()) {
                collectSessionIds(v, ids);
            }
        } else if (node instanceof List<?> list) {
            for (Object v : list) {
                collectSessionIds(v, ids);
            }
        }
    }

    private static String waitNewSessionId(Set<String> before, Duration timeout) throws InterruptedException {
        var deadline = System.currentTimeMillis() + timeout.toMillis();
        while (System.currentTimeMillis() < deadline) {
            Set<String> now = sessionIds();
            now.removeAll(before);
            if (!now.isEmpty()) {
                return now.iterator().next();
            }
            TimeUnit.MILLISECONDS.sleep(250);
        }
        return null;
    }

    private static String stripTrailingSlash(String url) {
        if (url == null || url.isEmpty()) {
            return url;
        }
        return url.endsWith("/") ? url.substring(0, url.length() - 1) : url;
    }
}
