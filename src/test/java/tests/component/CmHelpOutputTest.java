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
@DisplayName("CM help output fixture")
class CmHelpOutputTest {

    @Test
    @DisplayName("lists selenoid and selenoid-ui subcommands in help")
    void listsCoreSubcommandsInHelp() {
        var output = FixtureJson.load("fixtures/cm/help-output.txt");
        assertTrue(output.contains("selenoid"), () -> "Expected selenoid command in:\n" + output);
        assertTrue(output.contains("selenoid-ui"), () -> "Expected selenoid-ui command in:\n" + output);
        assertTrue(output.contains("version"), () -> "Expected version command in:\n" + output);
    }
}
