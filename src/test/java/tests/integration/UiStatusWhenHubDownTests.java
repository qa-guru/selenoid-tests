package tests.integration;

import annotations.Component;
import annotations.Layer;
import api.ui.UiStatusApi;
import helpers.StackHelper;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.Timeout;
import org.junit.jupiter.api.parallel.Execution;
import org.junit.jupiter.api.parallel.ExecutionMode;

import static io.qameta.allure.Allure.step;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertNotNull;

@Layer("integration")
@Component("selenoid-ui")
@Epic("selenoid-ui")
@Feature("UI status when hub down")
@DisplayName("UI status when hub down")
@Tag("integration")
@Tag("local-only")
@Execution(ExecutionMode.SAME_THREAD)
class UiStatusWhenHubDownTests {

    @Test
    @Tag("positive")
    @Timeout(60)
    @DisplayName("UI /status returns errors when hub is unavailable")
    void uiStatusReportsHubError() throws Exception {
        step("Stop hub", () -> {
            StackHelper.killHub();
            StackHelper.waitForHubDown(15_000);
        });

        try {
            var response = step("GET UI /status", UiStatusApi::fetchWhenHubUnavailable);

            step("Verify error payload", () -> {
                assertNotNull(response.errors());
                assertFalse(response.errors().isEmpty());
                assertNotNull(response.errors().getFirst().msg());
            });
        } finally {
            step("Restart hub", () -> {
                StackHelper.startHubDetached();
                StackHelper.waitForHubReady(30_000);
            });
        }
    }
}
