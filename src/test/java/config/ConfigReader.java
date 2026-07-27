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
     * <p>
     * Strips {@code user:pass@} from the URL: modern {@code fetch()} rejects relative requests
     * when the document URL includes credentials (breaks New Session → Create Session on prod).
     */
    public static String resolveUiBrowserUrl() {
        return resolveUiBrowserUrl(testConfig);
    }

    public static String resolveUiBrowserUrl(TestConfig config) {
        var remoteUrl = config.remoteUrl();
        if (remoteUrl == null || remoteUrl.isBlank()) {
            return stripUserInfo(stripTrailingSlash(resolveUiUrl(config)));
        }
        var uiUrl = config.uiUrl().trim();
        if (uiUrl.isEmpty()) {
            throw new IllegalStateException("Set uiUrl in config/${env}.properties");
        }
        var browserUrl = uiUrl
                .replace("127.0.0.1", "host.docker.internal")
                .replace("localhost", "host.docker.internal");
        return stripUserInfo(stripTrailingSlash(withSlash(browserUrl)));
    }

    /**
     * Basic-auth pair from {@code remoteUrl} / {@code apiBaseUrl} / {@code uiUrl}
     * ({@code https://user:pass@host/...}). Empty when absent.
     */
    public static String[] resolveHubBasicAuth() {
        return resolveHubBasicAuth(testConfig);
    }

    public static String[] resolveHubBasicAuth(TestConfig config) {
        for (var candidate : new String[]{config.remoteUrl(), config.apiBaseUrl(), config.uiUrl()}) {
            var creds = parseUserInfo(candidate);
            if (creds != null) {
                return creds;
            }
        }
        return new String[]{"", ""};
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

    /** Drop {@code user:pass@} so the browser document URL is fetch()-safe. */
    static String stripUserInfo(String url) {
        if (url == null || url.isBlank()) {
            return url;
        }
        var schemeSep = url.indexOf("://");
        if (schemeSep < 0) {
            return url;
        }
        var afterScheme = schemeSep + 3;
        var slash = url.indexOf('/', afterScheme);
        var authorityEnd = slash < 0 ? url.length() : slash;
        var authority = url.substring(afterScheme, authorityEnd);
        var at = authority.lastIndexOf('@');
        if (at < 0) {
            return url;
        }
        return url.substring(0, afterScheme) + authority.substring(at + 1) + url.substring(authorityEnd);
    }

    static String[] parseUserInfo(String url) {
        if (url == null || url.isBlank()) {
            return null;
        }
        var schemeSep = url.indexOf("://");
        if (schemeSep < 0) {
            return null;
        }
        var afterScheme = schemeSep + 3;
        var slash = url.indexOf('/', afterScheme);
        var authorityEnd = slash < 0 ? url.length() : slash;
        var authority = url.substring(afterScheme, authorityEnd);
        var at = authority.lastIndexOf('@');
        if (at <= 0) {
            return null;
        }
        var userInfo = authority.substring(0, at);
        var colon = userInfo.indexOf(':');
        if (colon <= 0) {
            return null;
        }
        return new String[]{userInfo.substring(0, colon), userInfo.substring(colon + 1)};
    }
}
