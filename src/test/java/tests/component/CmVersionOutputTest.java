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
@DisplayName("CM version output fixture")
class CmVersionOutputTest {

    @Test
    @DisplayName("parses git revision from cm version output")
    void parsesGitRevisionFromVersionOutput() {
        var output = FixtureJson.load("fixtures/cm/version-output.txt");
        assertTrue(output.contains("Git Revision:"), () -> "Expected revision line in:\n" + output);
        assertTrue(output.contains("UTC Build Time:"), () -> "Expected build time line in:\n" + output);
    }
}
