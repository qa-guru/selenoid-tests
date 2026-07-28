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
 *
 * <p>Two content depths (ADR 009): default {@code harContent=meta} (omit text) and opt-in
 * {@code harContent=bodies} (best-effort {@code content.text}; not ≡ Playwright {@code recordHar}).
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

    /**
     * Best-effort bodies gate on the local fixture: at least one HTTP entry with
     * {@code content.text} and matching {@code content.size} (hub ≥ v3.0.5 sets size from
     * decoded body). Explicitly <em>not</em> parity with Playwright {@code recordHar}.
     */
    private static final int BODIES_MIN_WITH_CONTENT_TEXT = 1;
    private static final int BODIES_MIN_WITH_CONTENT_SIZE = 1;

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

        HarStats wdHub = step("4) WebDriver → Selenoid hub enableHAR (meta default) → " + remoteUrl, () ->
                captureWebDriverHubHar(remoteUrl, null, "4-wd-hub-enableHAR"));
        stats.add(wdHub);

        HarStats pwHub = step("5) Playwright → Selenoid hub enableHAR meta (no client recordHar) → " + remoteUrl, () ->
                capturePlaywrightHubHar(remoteUrl, null, "5-pw-hub-enableHAR", gaps));
        if (pwHub != null) {
            stats.add(pwHub);
        }

        HarStats wdHubBodies = step("4b) WebDriver → hub enableHAR + harContent=bodies → " + remoteUrl, () ->
                captureWebDriverHubHar(remoteUrl, "bodies", "4b-wd-hub-enableHAR-bodies"));
        stats.add(wdHubBodies);

        HarStats pwHubBodies = step("5b) Playwright → hub enableHAR + harContent=bodies → " + remoteUrl, () ->
                capturePlaywrightHubHar(remoteUrl, "bodies", "5b-pw-hub-enableHAR-bodies", gaps));
        if (pwHubBodies != null) {
            stats.add(pwHubBodies);
        }

        String table = HarStats.formatTable(stats);
        System.out.println("\n=== Hub HAR completeness ===\n" + table + "\n");
        System.out.printf(
                "URL coverage vs baseline: wd-hub(meta)=%.0f%%%n",
                wdHub.urlCoverageOf(baseline) * 100);
        if (pwHub != null) {
            System.out.printf("URL coverage vs baseline: pw-hub(meta)=%.0f%%%n", pwHub.urlCoverageOf(baseline) * 100);
        }
        System.out.printf(
                "URL coverage vs baseline: wd-hub(bodies)=%.0f%%  withContentText=%d%n",
                wdHubBodies.urlCoverageOf(baseline) * 100,
                wdHubBodies.withContentText);
        if (pwHubBodies != null) {
            System.out.printf(
                    "URL coverage vs baseline: pw-hub(bodies)=%.0f%%  withContentText=%d%n",
                    pwHubBodies.urlCoverageOf(baseline) * 100,
                    pwHubBodies.withContentText);
        }
        if (!gaps.isEmpty()) {
            System.out.println("Gaps:\n- " + String.join("\n- ", gaps));
        }

        // Document known field gaps vs client recordHar (ADR 009 — not absolute omit).
        gaps.add("default hub enableHAR (harContent=meta|omit) omits content.text; bodies is opt-in best-effort");
        gaps.add("hub meta path may have content.size=0 and partial status for in-flight requests at quit");
        gaps.add("hub harContent=bodies (≥ v3.0.5) sets content.size from decoded body; still ≠ recordHar parity");
        gaps.add("do not combine hub enableHAR with client recordHar/HarCapture on one session");

        Map<String, Object> contentTextNote = new HashMap<>();
        contentTextNote.put(
                "meta",
                "default enableHAR / harContent=meta|omit: withContentText==0 on fixture");
        contentTextNote.put(
                "bodies",
                "enableHAR + harContent=bodies (hub ≥ v3.0.5): withContentText>="
                        + BODIES_MIN_WITH_CONTENT_TEXT
                        + ", withContentSize>="
                        + BODIES_MIN_WITH_CONTENT_SIZE
                        + " best-effort on fixture; not ≡ recordHar");

        Map<String, Object> summary = new HashMap<>();
        summary.put("targetLocal", localUrl);
        summary.put("targetRemote", remoteUrl);
        summary.put("rows", stats.stream().map(HarStats::toRow).toList());
        summary.put("wdHubUrlCoverage", wdHub.urlCoverageOf(baseline));
        summary.put("wdHubBodiesUrlCoverage", wdHubBodies.urlCoverageOf(baseline));
        summary.put("wdHubBodiesWithContentText", wdHubBodies.withContentText);
        if (pwHub != null) {
            summary.put("pwHubUrlCoverage", pwHub.urlCoverageOf(baseline));
        } else {
            summary.put("pwHubUrlCoverage", 0.0);
            summary.put("pwHubSkipped", true);
        }
        if (pwHubBodies != null) {
            summary.put("pwHubBodiesUrlCoverage", pwHubBodies.urlCoverageOf(baseline));
            summary.put("pwHubBodiesWithContentText", pwHubBodies.withContentText);
        } else {
            summary.put("pwHubBodiesSkipped", true);
        }
        summary.put("gaps", gaps);
        summary.put("contentTextNote", contentTextNote);
        summary.put("bodiesMinWithContentText", BODIES_MIN_WITH_CONTENT_TEXT);
        summary.put("bodiesMinWithContentSize", BODIES_MIN_WITH_CONTENT_SIZE);
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
        // Default meta path: content.text must stay empty. status/size may be 0 for in-flight
        // entries cancelled at session quit — documented in hub-summary gaps, not a hard fail.
        assertTrue(wdHub.withContentText == 0, "hub HAR meta default must not include content.text");
        assertTrue(wdHub.withRequestHeaders >= (int) Math.floor(wdHub.httpEntries * 0.80));

        assertTrue(pwHub != null, "Playwright hub-HAR must be produced (Chromium image :7070)");
        assertTrue(pwHub.withContentText == 0, "Playwright hub HAR meta default must not include content.text");
        assertTrue(
                pwHub.urlCoverageOf(baseline) >= 0.80,
                () -> "Playwright hub-HAR URL coverage < 80% of local baseline: "
                        + pwHub.urlCoverageOf(baseline)
                        + "\nbaseline=" + baseline.urls + "\ncandidate=" + pwHub.urls);

        // Bodies opt-in: best-effort threshold only (not recordHar parity). URL cov gate unchanged.
        assertTrue(
                wdHubBodies.urlCoverageOf(baseline) >= 0.80,
                () -> "WebDriver hub-HAR bodies URL coverage < 80% of local baseline: "
                        + wdHubBodies.urlCoverageOf(baseline));
        assertTrue(
                wdHubBodies.withContentText >= BODIES_MIN_WITH_CONTENT_TEXT,
                () -> "hub harContent=bodies expected withContentText>="
                        + BODIES_MIN_WITH_CONTENT_TEXT
                        + " on fixture (best-effort), got " + wdHubBodies.withContentText);
        assertTrue(
                wdHubBodies.withContentSize >= BODIES_MIN_WITH_CONTENT_SIZE,
                () -> "hub harContent=bodies expected withContentSize>="
                        + BODIES_MIN_WITH_CONTENT_SIZE
                        + " on fixture (hub ≥ v3.0.5), got " + wdHubBodies.withContentSize);
        assertTrue(pwHubBodies != null, "Playwright hub-HAR bodies must be produced");
        assertTrue(
                pwHubBodies.urlCoverageOf(baseline) >= 0.80,
                () -> "Playwright hub-HAR bodies URL coverage < 80% of local baseline: "
                        + pwHubBodies.urlCoverageOf(baseline));
        assertTrue(
                pwHubBodies.withContentText >= BODIES_MIN_WITH_CONTENT_TEXT,
                () -> "PW hub harContent=bodies expected withContentText>="
                        + BODIES_MIN_WITH_CONTENT_TEXT
                        + " on fixture (best-effort), got " + pwHubBodies.withContentText);
        assertTrue(
                pwHubBodies.withContentSize >= BODIES_MIN_WITH_CONTENT_SIZE,
                () -> "PW hub harContent=bodies expected withContentSize>="
                        + BODIES_MIN_WITH_CONTENT_SIZE
                        + " on fixture (hub ≥ v3.0.5), got " + pwHubBodies.withContentSize);
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

    /**
     * @param harContent {@code null}/{@code "meta"} = default meta path; {@code "bodies"} = opt-in
     * @param label      artifact stem under {@code build/har-compare/} (also HarStats label)
     */
    private static HarStats captureWebDriverHubHar(String url, String harContent, String label)
            throws Exception {
        Path harPath = OUT.resolve(label + ".har");
        Map<String, Object> selenoid = new HashMap<>();
        selenoid.put("enableVNC", true);
        selenoid.put("enableVideo", false);
        selenoid.put("enableHAR", true);
        selenoid.put("enableLog", false);
        selenoid.put("sessionTimeout", "2m");
        selenoid.put("name", label);
        if (harContent != null && !harContent.isBlank()) {
            selenoid.put("harContent", harContent);
        }

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
        return HarStats.fromBytes(label, harPath, body);
    }

    /**
     * Hub enableHAR on Playwright WS — no client recordHar.
     * Requires Chromium-family images with DevTools proxy on :7070.
     *
     * @param harContent {@code null} = meta default; {@code "bodies"} appends {@code harContent=} query
     */
    private static HarStats capturePlaywrightHubHar(
            String url, String harContent, String label, List<String> gaps) throws Exception {
        Path harPath = OUT.resolve(label + ".har");
        String base = ConfigReader.resolvePlaywrightWsEndpoint();
        String sep = base.contains("?") ? "&" : "?";
        StringBuilder qs = new StringBuilder("enableHAR=true&enableVideo=false&enableVNC=false&name=")
                .append(label);
        if (harContent != null && !harContent.isBlank()) {
            qs.append("&harContent=").append(harContent);
        }
        String ws = base + sep + qs;

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
            gaps.add("Playwright hub session id not observed in /status — cannot fetch /har (" + label + ")");
            return null;
        }

        String file = HubHarApi.waitForSessionHar(sessionId, Duration.ofSeconds(20));
        if (file == null) {
            gaps.add("Playwright hub-HAR missing for " + sessionId + " (" + label + ")"
                    + " — expected DevTools :7070 on playwright-chromium/chrome/msedge"
                    + " (HAR_CAPTURE_FAILED)");
            return null;
        }
        byte[] body = HubHarApi.download(file);
        Files.write(harPath, body);
        return HarStats.fromBytes(label, harPath, body);
    }

    private static Set<String> sessionIds() {
        Set<String> ids = new HashSet<>();
        collectSessionIds(HubStatusApi.fetch().browsers(), ids);
        return ids;
    }

    @SuppressWarnings("unchecked")
    private static void collectSessionIds(Object node, Set<String> ids) {
        if (node instanceof Map<?, ?> map) {
            Object sessionList = map.get("sessions");
            if (sessionList instanceof List<?> list) {
                for (Object item : list) {
                    if (item instanceof Map<?, ?> sess) {
                        Object id = sess.get("id");
                        if (id instanceof String s && !s.isBlank()) {
                            ids.add(s);
                        }
                    }
                }
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
