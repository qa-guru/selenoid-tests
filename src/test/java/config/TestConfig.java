package config;

import org.aeonbits.owner.Config;

@Config.LoadPolicy(Config.LoadType.MERGE)
@Config.Sources({
        "system:properties",
        "classpath:config/${env}.properties",
        "classpath:config/default.properties",
})
public interface TestConfig extends Config {

    @Key("attachBrowserConsoleLogs")
    @DefaultValue("false")
    boolean attachBrowserConsoleLogs();

    @Key("attachLastScreenshot")
    @DefaultValue("false")
    boolean attachLastScreenshot();

    @Key("attachPageSource")
    @DefaultValue("false")
    boolean attachPageSource();

    @Key("allureReportMode")
    @DefaultValue("allure3")
    String allureReportMode();

    @Key("logToConsole")
    @DefaultValue("true")
    boolean logToConsole();

    @Key("selenideLogToConsole")
    @DefaultValue("true")
    boolean selenideLogToConsole();

    @Key("rootLogLevel")
    @DefaultValue("info")
    String rootLogLevel();

    @Key("apiBaseUrl")
    @DefaultValue("")
    String apiBaseUrl();

    @Key("hubUrl")
    @DefaultValue("http://127.0.0.1:4444/")
    String hubUrl();

    @Key("uiUrl")
    @DefaultValue("http://127.0.0.1:8080/")
    String uiUrl();

    @Key("browser")
    @DefaultValue("chrome")
    String browser();

    @Key("browserSize")
    @DefaultValue("1920x1080")
    String browserSize();

    @Key("browserVersion")
    @DefaultValue("148.0")
    String browserVersion();

    @Key("closeBrowserAfterEach")
    @DefaultValue("true")
    boolean closeBrowserAfterEach();

    @Key("enableVnc")
    @DefaultValue("false")
    boolean enableVnc();

    @Key("enableVideo")
    @DefaultValue("false")
    boolean enableVideo();

    @Key("headless")
    @DefaultValue("true")
    boolean headless();

    @Key("remoteUrl")
    @DefaultValue("http://127.0.0.1:4444/wd/hub")
    String remoteUrl();

    @Key("smokeUrl")
    @DefaultValue("https://example.com/")
    String smokeUrl();

    @Key("playwrightWsEndpoint")
    @DefaultValue("ws://127.0.0.1:4444/playwright/playwright-chromium/1.61.1")
    String playwrightWsEndpoint();

    @Key("playwrightSessionName")
    @DefaultValue("java-playwright-tests")
    String playwrightSessionName();

    @Key("playwrightSessionTimeout")
    @DefaultValue("5m")
    String playwrightSessionTimeout();

    @Key("playwrightEnableVnc")
    @DefaultValue("false")
    boolean playwrightEnableVnc();

    @Key("playwrightEnableVideo")
    @DefaultValue("false")
    boolean playwrightEnableVideo();
}
