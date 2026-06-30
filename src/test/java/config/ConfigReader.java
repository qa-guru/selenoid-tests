package config;

import org.aeonbits.owner.ConfigFactory;

import java.net.URLEncoder;
import java.nio.charset.StandardCharsets;
import java.util.LinkedHashMap;
import java.util.Map;
import java.util.stream.Collectors;

public final class ConfigReader {
    public static final TestConfig testConfig = ConfigFactory.create(TestConfig.class);

    private ConfigReader() {
    }

    public static String resolveHubUrl() {
        return resolveHubUrl(testConfig);
    }

    public static String resolveHubUrl(TestConfig config) {
        var url = config.hubUrl().trim();
        if (url.isEmpty()) {
            throw new IllegalStateException("Set hubUrl in config/${env}.properties");
        }
        return withSlash(url);
    }

    public static String resolveUiUrl() {
        return resolveUiUrl(testConfig);
    }

    public static String resolveUiUrl(TestConfig config) {
        var url = config.uiUrl().trim();
        if (url.isEmpty()) {
            throw new IllegalStateException("Set uiUrl in config/${env}.properties");
        }
        return withSlash(url);
    }

    public static String resolvePlaywrightWsEndpoint() {
        return resolvePlaywrightWsEndpoint(testConfig);
    }

    public static String resolvePlaywrightWsEndpoint(TestConfig config) {
        var base = config.playwrightWsEndpoint().trim();
        if (base.isEmpty()) {
            throw new IllegalStateException("Set playwrightWsEndpoint in config/${env}.properties");
        }
        if (base.contains("?")) {
            return base;
        }

        Map<String, String> params = new LinkedHashMap<>();
        params.put("name", config.playwrightSessionName());
        params.put("sessionTimeout", config.playwrightSessionTimeout());
        params.put("enableVNC", String.valueOf(config.playwrightEnableVnc()));
        params.put("enableVideo", String.valueOf(config.playwrightEnableVideo()));

        var query = params.entrySet().stream()
                .map(entry -> encode(entry.getKey()) + "=" + encode(entry.getValue()))
                .collect(Collectors.joining("&"));
        return base + "?" + query;
    }

    private static String encode(String value) {
        return URLEncoder.encode(value, StandardCharsets.UTF_8);
    }

    private static String withSlash(String value) {
        return value.endsWith("/") ? value : value + "/";
    }
}
