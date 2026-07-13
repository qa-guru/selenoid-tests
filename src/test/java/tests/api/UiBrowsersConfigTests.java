package tests.api;

import annotations.Component;
import annotations.Layer;
import api.UiApiTestBase;
import api.ui.UiBrowsersConfigApi;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import io.qameta.allure.Story;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import static io.qameta.allure.Allure.step;
import static org.junit.jupiter.api.Assertions.assertNotNull;

@Layer("api")
@Component("selenoid-ui")
@Epic("selenoid-ui")
@Feature("UI browsers config")
@Story("UI browsers config API")
@DisplayName("UI browsers config API")
class UiBrowsersConfigTests extends UiApiTestBase {

    @Test
    @Tag("api")
    @Tag("positive")
    @DisplayName("GET /browsers-config returns browser catalog JSON")
    void browsersConfigReturnsCatalog() {
        var catalog = step("GET /browsers-config", UiBrowsersConfigApi::fetch);
        step("Verify catalog map is present", () -> assertNotNull(catalog));
    }
}
