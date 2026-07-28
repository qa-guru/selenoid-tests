package tests;

import allure.Attachments;
import annotations.Component;
import annotations.Layer;
import com.codeborne.selenide.Configuration;
import com.codeborne.selenide.Selenide;
import com.codeborne.selenide.WebDriverRunner;
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
import java.util.HashMap;
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
 * Prod smoke: client {@link HarCapture} {@code HarContentMode.BODIES} on warm Chrome via
 * WebDriver — one writer, no hub {@code enableHAR}.
 *
 * <p>Short target ({@code smokeUrl}, default example.com). Gates: {@code withContentText >= 1}
 * (required); {@code withContentSize >= 1} best-effort (logged, not a hard fail).
 */
@Layer("e2e")
@Component("selenoid")
@Epic("selenoid")
@Feature("HAR")
@Story("HarCapture BODIES on prod Selenoid")
@DisplayName("HarCapture prod e2e")
@Tag("smoke")
@Tag("har-capture")
@Tag("positive")
class HarCaptureProdE2eTests {

    private static final TestConfig config = ConfigReader.testConfig;
    private static final Json JSON = new Json();
    private static final Path OUT = Path.of("build/har-compare");
    private static final Path PROD_STEP5 = OUT.resolve("prod-step5");
    private static final int BODIES_MIN_WITH_CONTENT_TEXT = 1;

    @Test
    @DisplayName("Selenide HarCapture BODIES → valid HAR with content.text on prod")
    void harCaptureBodiesProducesTextOnProd() throws Exception {
        Files.createDirectories(OUT);
        Files.createDirectories(PROD_STEP5);

        String targetUrl = stripTrailingSlash(config.smokeUrl());
        String label = "3b-selenide-HarCapture-bodies-prod";
        Path harPath = OUT.resolve("3b-selenide-HarCapture-bodies.har");
        Path prodHarPath = OUT.resolve("3b-selenide-HarCapture-bodies-prod.har");
        Path step5HarPath = PROD_STEP5.resolve("3b-selenide-HarCapture-bodies-prod.har");

        HarStats stats = step("Capture HarCapture BODIES on " + targetUrl, () ->
                captureSelenideHarCaptureBodies(targetUrl, label, harPath));

        Files.copy(harPath, prodHarPath, java.nio.file.StandardCopyOption.REPLACE_EXISTING);
        Files.copy(harPath, step5HarPath, java.nio.file.StandardCopyOption.REPLACE_EXISTING);
        Attachments.har(prodHarPath);

        boolean sizeGateOk = stats.withContentSize >= 1;
        if (!sizeGateOk) {
            System.out.println(
                    "NOTE: HarCapture prod bodies withContentSize="
                            + stats.withContentSize
                            + " (best-effort gate; text gate is authoritative)");
        }

        Map<String, Object> summary = new HashMap<>();
        summary.put("base", stripTrailingSlash(config.hubUrl()));
        summary.put("nav", targetUrl);
        summary.put("writer", "harcapture-selenide");
        summary.put("mode", "bodies");
        summary.put("label", label);
        summary.put("bodiesMinWithContentText", BODIES_MIN_WITH_CONTENT_TEXT);
        summary.put("stats", stats.toRow());
        summary.put("gates", Map.of(
                "withContentText", stats.withContentText >= BODIES_MIN_WITH_CONTENT_TEXT,
                "withContentSizeBestEffort", sizeGateOk));
        summary.put("artifacts", Map.of(
                "har", prodHarPath.toString(),
                "prodStep5", step5HarPath.toString()));
        summary.put(
                "note",
                "One writer (HarCapture only, enableHAR=false); not ≡ recordHar; "
                        + "size gate best-effort on prod");

        Files.writeString(
                OUT.resolve("harcapture-prod-summary.json"),
                JSON.toJson(summary),
                StandardCharsets.UTF_8);
        Files.writeString(
                PROD_STEP5.resolve("harcapture-prod-summary.json"),
                JSON.toJson(summary),
                StandardCharsets.UTF_8);

        step("Assert HarCapture prod bodies gates", () -> {
            assertTrue(stats.httpEntries >= 1, "expected ≥1 http entry on prod smoke");
            assertTrue(
                    stats.withContentText >= BODIES_MIN_WITH_CONTENT_TEXT,
                    () -> "HarCapture BODIES prod expected withContentText>="
                            + BODIES_MIN_WITH_CONTENT_TEXT
                            + ", got "
                            + stats.withContentText);
        });
    }

    private static HarStats captureSelenideHarCaptureBodies(String url, String label, Path harPath)
            throws Exception {
        Configuration.browser = "chrome";
        Configuration.browserVersion = config.chromeVersion();
        Configuration.browserSize = config.browserSize();
        Configuration.headless = config.headless();
        Configuration.remote = config.remoteUrl();
        Configuration.timeout = 30_000;

        ChromeOptions chrome = new ChromeOptions();
        chrome.addArguments(
                "--headless=new",
                "--no-sandbox",
                "--disable-dev-shm-usage",
                "--disable-features=ServiceWorker");
        HarCapture.enablePerformanceLogging(chrome);

        MutableCapabilities caps = new MutableCapabilities();
        caps.setCapability(ChromeOptions.CAPABILITY, chrome);
        Map<String, Object> selenoid = new HashMap<>();
        selenoid.put("enableVNC", false);
        selenoid.put("enableVideo", false);
        selenoid.put("enableHAR", false);
        selenoid.put("headless", config.headless());
        selenoid.put("name", label);
        caps.setCapability("selenoid:options", selenoid);
        HarCapture.enablePerformanceLogging(caps);
        Configuration.browserCapabilities = caps;

        try {
            open(url);
            Selenide.sleep(2000);
            byte[] har = HarCapture.collectHarJson(HarCapture.HarContentMode.BODIES)
                    .orElseThrow(() -> new IllegalStateException("HarCapture BODIES produced empty HAR"));
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
