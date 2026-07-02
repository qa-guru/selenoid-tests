package tests.integration;

import annotations.Component;
import annotations.Layer;
import helpers.CmCli;
import io.qameta.allure.Epic;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import static io.qameta.allure.Allure.step;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertTrue;

@Layer("integration")
@Component("cm")
@Epic("cm")
@DisplayName("CM CLI version")
class CmCliVersionTests {

    @Test
    @Tag("integration")
    @Tag("cm")
    @Tag("local-only")
    @DisplayName("cm version exits zero and prints revision")
    void cmVersionPrintsRevision() {
        var result = step("Run cm version", CmCli::version);

        step("Verify exit code and output", () -> {
            assertEquals(0, result.exitCode(), () -> "cm version failed:\n" + result.output());
            assertTrue(result.output().contains("Git Revision:"), () -> "Unexpected output:\n" + result.output());
        });
    }

    @Test
    @Tag("integration")
    @Tag("cm")
    @Tag("local-only")
    @DisplayName("cm --help exits zero and lists subcommands")
    void cmHelpListsSubcommands() {
        var result = step("Run cm --help", CmCli::help);

        step("Verify exit code and output", () -> {
            assertEquals(0, result.exitCode(), () -> "cm --help failed:\n" + result.output());
            assertTrue(result.output().contains("selenoid-ui"), () -> "Unexpected output:\n" + result.output());
        });
    }
}
