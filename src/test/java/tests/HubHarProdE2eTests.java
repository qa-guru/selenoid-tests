package tests;

import allure.Attachments;
import annotations.Component;
import annotations.Layer;
import api.hub.HubHarApi;
import api.hub.HubSessionApi;
import api.hub.HubStatusApi;
import com.microsoft.playwright.Browser;
import com.microsoft.playwright.BrowserContext;
import com.microsoft.playwright.Page;
import com.microsoft.playwright.Playwright;
import com.microsoft.playwright.options.ServiceWorkerPolicy;
import config.ConfigReader;
import config.TestConfig;
import helpers.HarStats;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import io.qameta.allure.Story;
import java.net.URI;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.time.Duration;
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
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * Prod smoke: hub {@code enableHAR} → {@code /har/<session>.har} is valid HAR 1.2.
 *
 * <p>One writer per session: hub only — no client {@code recordHar}/HarCapture.
 */
@Layer("e2e")
@Component("selenoid")
@Epic("selenoid")
@Feature("HAR")
@Story("Hub enableHAR produces downloadable /har artifact")
@DisplayName("Hub HAR prod e2e")
@Tag("smoke")
@Tag("hub-har")
@Tag("positive")
class HubHarProdE2eTests {

    private static final TestConfig config = ConfigReader.testConfig;
    private static final Json JSON = new Json();
    private static final Path OUT = Path.of("build/har-prod-e2e");

    @Test
    @DisplayName("WebDriver enableHAR → /har is valid HAR 1.2 with target URL")
    void webDriverEnableHarProducesValidArtifact() throws Exception {
        Files.createDirectories(OUT);
        String targetUrl = stripTrailingSlash(config.smokeUrl());
        String targetHost = hostOf(targetUrl);

        Map<String, Object> selenoid = new HashMap<>();
        selenoid.put("enableVNC", false);
        selenoid.put("enableVideo", false);
        selenoid.put("enableHAR", true);
        selenoid.put("enableLog", false);
        selenoid.put("sessionTimeout", "2m");
        selenoid.put("name", "hub-har-wd-prod-e2e");

        String sessionId = step("Create WebDriver session with enableHAR", () ->
                HubSessionApi.createWithSelenoidOptions(config.chromeVersion(), selenoid));
        try {
            step("Navigate to " + targetUrl, () -> {
                HubSessionApi.navigate(sessionId, targetUrl);
                TimeUnit.MILLISECONDS.sleep(2000);
            });
        } finally {
            step("Delete WebDriver session", () -> HubSessionApi.delete(sessionId));
        }

        byte[] body = step("Wait and download /har for " + sessionId, () -> {
            String file = HubHarApi.waitForSessionHar(sessionId, Duration.ofSeconds(30));
            assertNotNull(file, () -> "expected hub HAR for WD session " + sessionId);
            return HubHarApi.download(file);
        });

        Path harPath = OUT.resolve("wd-" + sessionId + ".har");
        Files.write(harPath, body);
        Attachments.har(harPath);

        step("Assert HAR 1.2 + entries + target host", () ->
                assertValidHubHar(body, harPath, "wd-hub-enableHAR", targetHost));
    }

    @Test
    @DisplayName("Playwright enableHAR → /har is valid HAR 1.2 with target URL")
    void playwrightEnableHarProducesValidArtifact() throws Exception {
        Files.createDirectories(OUT);
        String targetUrl = stripTrailingSlash(config.smokeUrl());
        String targetHost = hostOf(targetUrl);

        String base = ConfigReader.resolvePlaywrightWsEndpoint();
        String sep = base.contains("?") ? "&" : "?";
        String ws = base + sep + "enableHAR=true&enableVideo=false&name=hub-har-pw-prod-e2e";

        Set<String> before = sessionIds();
        String sessionId = step("Connect Playwright with enableHAR and navigate", () -> {
            String id;
            try (Playwright pw = Playwright.create()) {
                Browser browser = pw.chromium().connect(ws);
                id = waitNewSessionId(before, Duration.ofSeconds(30));
                BrowserContext ctx = browser.newContext(new Browser.NewContextOptions()
                        .setServiceWorkers(ServiceWorkerPolicy.BLOCK));
                Page page = ctx.newPage();
                // Hub attaches to /page asynchronously after newPage(); yield so
                // Network.enable wins the race before the first navigation.
                page.waitForTimeout(750);
                page.navigate(targetUrl);
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
            return id;
        });

        assertNotNull(sessionId, "Playwright hub session id not observed in /status");

        byte[] body = step("Wait and download /har for " + sessionId, () -> {
            String file = HubHarApi.waitForSessionHar(sessionId, Duration.ofSeconds(30));
            assertNotNull(file, () -> "expected hub HAR for PW session " + sessionId
                    + " — Chromium-family image with DevTools :7070");
            return HubHarApi.download(file);
        });

        Path harPath = OUT.resolve("pw-" + sessionId + ".har");
        Files.write(harPath, body);
        Attachments.har(harPath);

        step("Assert HAR 1.2 + entries + target host", () ->
                assertValidHubHar(body, harPath, "pw-hub-enableHAR", targetHost));
    }

    @SuppressWarnings("unchecked")
    private static void assertValidHubHar(byte[] body, Path harPath, String label, String targetHost)
            throws Exception {
        Map<String, Object> root = JSON.toType(new String(body, StandardCharsets.UTF_8), Map.class);
        Map<String, Object> log = (Map<String, Object>) root.get("log");
        assertNotNull(log, "HAR must contain log");
        assertEquals("1.2", String.valueOf(log.get("version")), "HAR log.version");

        List<?> entries = (List<?>) log.get("entries");
        assertNotNull(entries, "HAR must contain log.entries");
        assertFalse(entries.isEmpty(), "HAR entries must not be empty");

        boolean hasTarget = false;
        boolean hasMethod = false;
        for (Object raw : entries) {
            if (!(raw instanceof Map<?, ?> entryMap)) {
                continue;
            }
            @SuppressWarnings("unchecked")
            Map<String, Object> entry = (Map<String, Object>) entryMap;
            Object reqObj = entry.get("request");
            Map<String, Object> req = reqObj instanceof Map<?, ?> rm
                    ? (Map<String, Object>) rm
                    : Map.of();
            String url = String.valueOf(req.getOrDefault("url", ""));
            if (url.contains(targetHost)) {
                hasTarget = true;
            }
            Object method = req.get("method");
            if (method instanceof String m && !m.isBlank()) {
                hasMethod = true;
            }
            Object time = entry.get("time");
            if (time instanceof Number n) {
                assertTrue(n.doubleValue() >= 0, "entry.time must be >= 0 when present");
            }
        }
        assertTrue(hasTarget, () -> "expected ≥1 request.url containing host " + targetHost);
        assertTrue(hasMethod, "expected ≥1 request.method");

        HarStats stats = HarStats.fromBytes(label, harPath, body);
        assertTrue(stats.entries > 0, "HarStats.entries must be > 0");
        assertTrue(stats.httpEntries > 0, "HarStats.httpEntries must be > 0");
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

    private static String hostOf(String url) {
        return URI.create(url).getHost();
    }

    private static String stripTrailingSlash(String url) {
        if (url == null || url.isEmpty()) {
            return url;
        }
        return url.endsWith("/") ? url.substring(0, url.length() - 1) : url;
    }
}
