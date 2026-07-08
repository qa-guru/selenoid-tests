package tests.integration;

import annotations.Component;
import annotations.Layer;
import api.hub.HubStatusApi;
import api.ui.UiStatusApi;
import config.ConfigReader;
import helpers.CmInstallerHelper;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.Timeout;
import org.junit.jupiter.api.parallel.Execution;
import org.junit.jupiter.api.parallel.ExecutionMode;

import java.io.IOException;
import java.nio.file.Files;

import static io.qameta.allure.Allure.step;
import static io.restassured.path.json.JsonPath.from;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertTrue;

@Layer("integration")
@Component("cm")
@Epic("cm")
@Feature("Installer lifecycle")
@DisplayName("CM installer lifecycle")
@Tag("integration")
@Tag("cm")
@Execution(ExecutionMode.SAME_THREAD)
@Timeout(300)
class CmInstallerLifecycleTests {

    private CmInstallerHelper installer;

    @BeforeEach
    void setUp() throws IOException {
        installer = CmInstallerHelper.withTempConfigDir();
        installer.stopAll();
    }

    @AfterEach
    void tearDown() {
        if (installer != null) {
            installer.stopAll();
            installer.deleteConfigDir();
        }
    }

    @Test
    @Tag("positive")
    @DisplayName("configure writes browsers.json from dev catalog")
    void configureWritesBrowsersJson() throws IOException {
        var result = step("cm selenoid configure -n", () -> {
            try {
                return installer.configure();
            } catch (IOException e) {
                throw new IllegalStateException("configure failed: " + e.getMessage(), e);
            }
        });
        result.requireSuccess("configure");

        step("Verify browsers.json exists and contains chrome", () -> {
            var configPath = installer.browsersJsonPath();
            assertTrue(configPath.toFile().isFile(), () -> "Missing " + configPath);

            var content = Files.readString(configPath);
            var chromeDefault = from(content).getString("chrome.default");
            assertFalse(chromeDefault.isBlank(), () -> "Expected chrome.default in " + configPath);
        });

        var status = step("cm selenoid status", installer::statusHub);
        status.requireSuccess("status");
        step("Verify status reports configuration file", () ->
                assertTrue(status.output().contains("configuration file"),
                        () -> "Unexpected status output:\n" + status.output()));
    }

    @Test
    @Tag("positive")
    @DisplayName("start exposes hub /status")
    void startHubExposesStatusEndpoint() throws Exception {
        step("Configure CM stack", () ->
                installer.configure().requireSuccess("configure"));
        step("Start hub via cm", () -> {
            var result = installer.startHub();
            result.requireSuccess("start hub");
        });

        step("Wait for hub readiness", () -> installer.waitForHubReady(120_000));

        var status = step("cm selenoid status", installer::statusHub);
        status.requireSuccess("status");
        step("Verify hub container is running", () ->
                assertTrue(status.output().contains("container is running"),
                        () -> "Unexpected status output:\n" + status.output()));

        step("GET hub /status", () ->
                HubStatusApi.fetchFrom(ConfigReader.resolveCmHubUrl()));
    }

    @Test
    @Tag("positive")
    @DisplayName("full stack start — UI /status mirrors hub counters")
    void startFullStackUiProxiesHub() throws Exception {
        step("Configure CM stack", () ->
                installer.configure().requireSuccess("configure"));
        step("Start hub via cm", () ->
                installer.startHub().requireSuccess("start hub"));
        step("Wait for hub readiness", () ->
                installer.waitForHubReady(120_000));
        step("Start UI via cm", () ->
                installer.startUi().requireSuccess("start UI"));
        step("Wait for UI readiness", () ->
                installer.waitForUiReady(90_000));

        var hubStatus = step("GET hub /status", () ->
                HubStatusApi.fetchFrom(ConfigReader.resolveCmHubUrl()));
        var uiStatus = step("GET UI /status", () ->
                UiStatusApi.fetchFrom(ConfigReader.resolveCmUiUrl()));

        step("Verify proxied counters match hub", () -> {
            assertEquals(hubStatus.total(), uiStatus.total());
            assertEquals(hubStatus.used(), uiStatus.used());
            assertEquals(hubStatus.queued(), uiStatus.queued());
            assertEquals(hubStatus.pending(), uiStatus.pending());
        });
    }
}
