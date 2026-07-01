package api;

import config.ConfigReader;
import config.TestConfig;
import io.qameta.allure.restassured.AllureRestAssured;
import io.restassured.RestAssured;
import org.junit.jupiter.api.BeforeAll;

public class UiApiTestBase {

    protected static final TestConfig config = ConfigReader.testConfig;

    @BeforeAll
    static void setupRestAssured() {
        RestAssured.baseURI = ConfigReader.resolveUiUrl();
        RestAssured.enableLoggingOfRequestAndResponseIfValidationFails();
        if (config.logToConsole()) {
            RestAssured.filters(new AllureRestAssured());
        }
    }
}
