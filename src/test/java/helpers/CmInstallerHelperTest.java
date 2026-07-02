package helpers;

import annotations.Layer;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import java.nio.file.Files;

import static org.junit.jupiter.api.Assertions.assertTrue;

@Layer("unit")
@DisplayName("CmInstallerHelper paths")
class CmInstallerHelperTest {

    @Test
    @DisplayName("withTempConfigDir creates isolated config directory")
    void withTempConfigDirCreatesConfigDirectory() throws Exception {
        var helper = CmInstallerHelper.withTempConfigDir();
        try {
            assertTrue(Files.isDirectory(helper.configDir()));
            assertTrue(Files.notExists(helper.browsersJsonPath()));
        } finally {
            helper.deleteConfigDir();
        }
    }
}
