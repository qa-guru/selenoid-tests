package tests.api;

import annotations.Component;
import annotations.Layer;
import api.ui.UiStatusApi;
import config.ConfigReader;
import io.qameta.allure.Epic;
import io.qameta.allure.Feature;
import io.qameta.allure.restassured.AllureRestAssured;
import io.restassured.RestAssured;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;

import static io.qameta.allure.Allure.step;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

@Layer("api")
@Component("cm")
@Epic("cm")
@Feature("CM-managed UI")
@DisplayName("CM UI status API")
class CmUiStatusApiTests {

    @BeforeAll
    static void setupRestAssured() {
        RestAssured.baseURI = ConfigReader.resolveCmUiUrl();
        RestAssured.enableLoggingOfRequestAndResponseIfValidationFails();
        if (ConfigReader.testConfig.logToConsole()) {
            RestAssured.filters(new AllureRestAssured());
        }
    }

    @Test
    @Tag("api")
    @Tag("cm")
    @Tag("positive")
    @DisplayName("GET /status on CM UI port returns proxied hub counters")
    void cmUiStatusReturnsProxiedCounters() {
        var status = step("GET CM UI /status", () ->
                UiStatusApi.fetchFrom(ConfigReader.resolveCmUiUrl()));

        step("Verify UI status counters are present", () -> {
            assertTrue(status.total() >= 0);
            assertTrue(status.used() >= 0);
            assertNotNull(status.browsers());
        });
    }
}
