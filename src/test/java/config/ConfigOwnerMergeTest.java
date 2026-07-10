package config;

import annotations.Component;
import annotations.Layer;
import org.aeonbits.owner.ConfigFactory;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;

@Component("selenoid")
@Layer("unit")
@DisplayName("Owner config merge")
class ConfigOwnerMergeTest {

    @Test
    @DisplayName("system property hubUrl overrides file default")
    void systemPropertyOverridesFileDefault() {
        var previous = System.getProperty("hubUrl");
        try {
            System.setProperty("hubUrl", "http://override.example.com:4444");
            var config = ConfigFactory.create(TestConfig.class);
            assertEquals("http://override.example.com:4444", config.hubUrl());
        } finally {
            if (previous == null) {
                System.clearProperty("hubUrl");
            } else {
                System.setProperty("hubUrl", previous);
            }
        }
    }
}
