package tests.integration;

import annotations.Component;
import annotations.Layer;
import api.hub.HubSessionApi;
import api.hub.HubStatusApi;
import config.ConfigReader;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import io.qameta.allure.Story;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.parallel.Execution;
import org.junit.jupiter.api.parallel.ExecutionMode;
import org.junit.jupiter.api.parallel.ResourceAccessMode;
import org.junit.jupiter.api.parallel.ResourceLock;

import java.util.Map;

import static io.qameta.allure.Allure.step;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;

@Layer("integration")
@Component("selenoid")
@Epic("selenoid")
@Feature("Hub status with session")
@Story("Hub status with session")
@DisplayName("Hub status browsers with active session")
@Execution(ExecutionMode.SAME_THREAD)
@ResourceLock(value = "hubSessions", mode = ResourceAccessMode.READ_WRITE)
class HubStatusSessionIntegrationTests {

    @Test
    @Tag("integration")
    @Tag("positive")
    @DisplayName("Hub /status browsers map reflects active session browser family")
    void statusBrowsersReflectActiveSession() {
        var sessionId = step("Create hub session", () -> HubSessionApi.create());
        try {
            var status = step("GET hub /status", HubStatusApi::fetch);
            step("Verify active session browser family is listed", () -> {
                assertEquals(1, status.used());
                assertNotNull(status.browsers());
                assertNotNull(status.browsers().get("chrome"));
                @SuppressWarnings("unchecked")
                var chromeVersions = (Map<String, Object>) status.browsers().get("chrome");
                assertNotNull(chromeVersions.get(ConfigReader.testConfig.browserVersion()));
            });
        } finally {
            step("Delete hub session", () -> HubSessionApi.delete(sessionId));
        }
    }
}
