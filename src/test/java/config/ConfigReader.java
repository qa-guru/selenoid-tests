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

    public static String resolveApiBaseUrl() {
        return resolveApiBaseUrl(testConfig);
    }

    public static String resolveApiBaseUrl(TestConfig config) {
        var apiUrl = config.apiBaseUrl().trim();
        if (!apiUrl.isEmpty()) {
            return withSlash(apiUrl);
        }

        var hubUrl = config.hubUrl().trim();
        if (!hubUrl.isEmpty()) {
            return withSlash(hubUrl);
        }

        throw new IllegalStateException("Set apiBaseUrl or hubUrl in config/${env}.properties");
    }

    public static String resolveHubStatusPath() {
        return resolveHubStatusPath(testConfig);
    }

    public static String resolveHubStatusPath(TestConfig config) {
        var path = config.hubStatusPath().trim();
        if (path.isEmpty()) {
            return "/status";
        }
        return path.startsWith("/") ? path : "/" + path;
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

    /**
     * Selenide baseUrl when the browser runs in a Selenoid container: loopback in uiUrl
     * is unreachable from inside Docker — use host.docker.internal (requires hosts in browsers.json).
     */
    public static String resolveUiBrowserUrl() {
        return resolveUiBrowserUrl(testConfig);
    }

    public static String resolveUiBrowserUrl(TestConfig config) {
        var remoteUrl = config.remoteUrl();
        if (remoteUrl == null || remoteUrl.isBlank()) {
            return stripTrailingSlash(resolveUiUrl(config));
        }
        var uiUrl = config.uiUrl().trim();
        if (uiUrl.isEmpty()) {
            throw new IllegalStateException("Set uiUrl in config/${env}.properties");
        }
        var browserUrl = uiUrl
                .replace("127.0.0.1", "host.docker.internal")
                .replace("localhost", "host.docker.internal");
        return stripTrailingSlash(withSlash(browserUrl));
    }

    public static String resolveCmHubUrl() {
        return resolveCmHubUrl(testConfig);
    }

    public static String resolveCmHubUrl(TestConfig config) {
        return withSlash("http://127.0.0.1:" + config.cmHubPort());
    }

    public static String resolveCmUiUrl() {
        return resolveCmUiUrl(testConfig);
    }

    public static String resolveCmUiUrl(TestConfig config) {
        return withSlash("http://127.0.0.1:" + config.cmUiPort());
    }

    public static String resolveCmRemoteUrl() {
        return resolveCmRemoteUrl(testConfig);
    }

    public static String resolveCmRemoteUrl(TestConfig config) {
        return resolveCmHubUrl(config) + "wd/hub";
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

    public static String resolvePlaywrightWsEndpoint(TestConfig config, String playwrightBrowser) {
        var endpoint = config.playwrightWsEndpoint().trim();
        var pathOnly = endpoint.contains("?") ? endpoint.substring(0, endpoint.indexOf('?')) : endpoint;
        var browserPath = pathOnly.replace("playwright-chromium", playwrightBrowser);
        if (endpoint.contains("?")) {
            return browserPath + endpoint.substring(endpoint.indexOf('?'));
        }

        Map<String, String> params = new LinkedHashMap<>();
        params.put("name", config.playwrightSessionName());
        params.put("sessionTimeout", config.playwrightSessionTimeout());
        params.put("enableVNC", String.valueOf(config.playwrightEnableVnc()));
        params.put("enableVideo", String.valueOf(config.playwrightEnableVideo()));

        var query = params.entrySet().stream()
                .map(entry -> encode(entry.getKey()) + "=" + encode(entry.getValue()))
                .collect(Collectors.joining("&"));
        return browserPath + "?" + query;
    }

    private static String encode(String value) {
        return URLEncoder.encode(value, StandardCharsets.UTF_8);
    }

    private static String withSlash(String value) {
        return value.endsWith("/") ? value : value + "/";
    }

    private static String stripTrailingSlash(String value) {
        return value.endsWith("/") ? value.substring(0, value.length() - 1) : value;
    }
}
