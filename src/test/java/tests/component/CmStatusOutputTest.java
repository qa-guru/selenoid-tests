package tests.component;

import annotations.Component;
import annotations.Layer;
import io.qameta.allure.Epic;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertTrue;

@Component("cm")
@Layer("component")
@Epic("cm")
@DisplayName("CM status output fixture")
class CmStatusOutputTest {

    @Test
    @DisplayName("detects running container in cm status output")
    void detectsRunningContainerInStatusOutput() {
        var output = FixtureJson.load("fixtures/cm/status-running.txt");
        assertTrue(output.contains("Configuration file:"), () -> "Expected configuration file line in:\n" + output);
        assertTrue(output.contains("container is running on port 4445"), () -> "Expected running container in:\n" + output);
    }
}
