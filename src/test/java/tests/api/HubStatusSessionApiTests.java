package tests.api;

import annotations.Component;
import annotations.Layer;
import api.ApiTestBase;
import api.hub.HubSessionApi;
import api.hub.HubStatusApi;
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

import java.time.Duration;

import static io.qameta.allure.Allure.step;
import static org.junit.jupiter.api.Assertions.assertEquals;

@Layer("api")
@Component("selenoid")
@Epic("selenoid")
@Feature("Hub status with session")
@Story("Hub status with session")
@DisplayName("Hub status session counters API")
@Execution(ExecutionMode.SAME_THREAD)
@ResourceLock(value = "hubSessions", mode = ResourceAccessMode.READ_WRITE)
class HubStatusSessionApiTests extends ApiTestBase {

    @Test
    @Tag("api")
    @Tag("positive")
    @DisplayName("GET /status used counter tracks session lifecycle")
    void statusUsedCounterTracksSessionLifecycle() {
        var usedBefore = step("Snapshot hub used counter", () -> HubStatusApi.fetch().used());
        var sessionId = step("Create hub session", () -> HubSessionApi.create(config));
        try {
            step("Verify used counter incremented", () ->
                    assertEquals(usedBefore + 1,
                            HubStatusApi.waitUntilUsed(usedBefore + 1, Duration.ofSeconds(30)).used()));
        } finally {
            step("Delete hub session", () -> HubSessionApi.delete(sessionId));
        }
        step("Verify used counter returned to baseline", () ->
                assertEquals(usedBefore,
                        HubStatusApi.waitUntilUsed(usedBefore, Duration.ofSeconds(30)).used()));
    }
}
