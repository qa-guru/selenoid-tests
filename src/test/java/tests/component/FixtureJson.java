package tests.component;

import io.restassured.path.json.JsonPath;

import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;

final class FixtureJson {

    private FixtureJson() {
    }

    static String loadProject(String relativeToProjectRoot) {
        var path = Path.of(System.getProperty("user.dir")).resolve(relativeToProjectRoot).normalize();
        try {
            return Files.readString(path, StandardCharsets.UTF_8);
        } catch (IOException e) {
            throw new IllegalStateException("Failed to read project fixture: " + path, e);
        }
    }

    static String load(String resourcePath) {
        try (var stream = FixtureJson.class.getClassLoader().getResourceAsStream(resourcePath)) {
            if (stream == null) {
                throw new IllegalStateException("Fixture not found: " + resourcePath);
            }
            return new String(stream.readAllBytes(), StandardCharsets.UTF_8);
        } catch (IOException e) {
            throw new IllegalStateException("Failed to read fixture: " + resourcePath, e);
        }
    }

    static <T> T parse(String resourcePath, Class<T> type) {
        return JsonPath.from(load(resourcePath)).getObject("", type);
    }
}
