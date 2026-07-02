package helpers;

import annotations.Component;
import annotations.Layer;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertDoesNotThrow;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

@Component("cm")
@Layer("unit")
@DisplayName("CmRunResult")
class CmRunResultTest {

    @Test
    @DisplayName("requireSuccess passes on exit code zero")
    void requireSuccessPassesOnZeroExit() {
        var result = new CmInstallerHelper.CmRunResult(0, "ok");
        assertDoesNotThrow(() -> result.requireSuccess("configure"));
    }

    @Test
    @DisplayName("requireSuccess fails on non-zero exit code")
    void requireSuccessFailsOnNonZeroExit() {
        var result = new CmInstallerHelper.CmRunResult(1, "boom");
        var error = assertThrows(IllegalStateException.class, () -> result.requireSuccess("start"));
        assertTrue(error.getMessage().contains("start"));
        assertTrue(error.getMessage().contains("boom"));
    }
}
