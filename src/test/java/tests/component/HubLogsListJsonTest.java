package tests.component;

import annotations.Component;
import annotations.Layer;
import io.qameta.allure.Epic;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertTrue;

@Component("selenoid")
@Layer("component")
@Epic("selenoid")
@DisplayName("Hub logs list fixture")
class HubLogsListJsonTest {

    @Test
    @DisplayName("parses hub logs JSON array with session file names")
    void parsesHubLogsJsonArray() {
        var files = io.restassured.path.json.JsonPath.from(
                FixtureJson.load("fixtures/hub/logs-list.json")).getList("");
        assertFalse(files.isEmpty());
        assertTrue(files.stream().map(String::valueOf).anyMatch(name -> name.endsWith(".log")));
    }
}
