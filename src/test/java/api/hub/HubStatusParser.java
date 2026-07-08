package api.hub;

import io.restassured.path.json.JsonPath;
import io.restassured.response.ResponseBodyExtractionOptions;

/**
 * Parses hub capacity JSON: flat hub {@code /status} or UI-shaped {@code /status} with {@code .state}.
 */
public final class HubStatusParser {

    private HubStatusParser() {
    }

    public static HubStatus parse(ResponseBodyExtractionOptions body) {
        return parseJsonPath(body.jsonPath());
    }

    public static HubStatus parseJson(String json) {
        return parseJsonPath(JsonPath.from(json));
    }

    private static HubStatus parseJsonPath(JsonPath jsonPath) {
        if (jsonPath.get("state") != null) {
            return jsonPath.getObject("state", HubStatus.class);
        }
        return jsonPath.getObject("", HubStatus.class);
    }
}
