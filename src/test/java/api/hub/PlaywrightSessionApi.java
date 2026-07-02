package api.hub;

import com.microsoft.playwright.Browser;
import com.microsoft.playwright.Playwright;
import config.ConfigReader;
import io.qameta.allure.Step;

public final class PlaywrightSessionApi {

    private PlaywrightSessionApi() {
    }

    @Step("Connect Playwright browser via hub WS")
    public static Browser connect(Playwright playwright) {
        return connect(playwright, ConfigReader.resolvePlaywrightWsEndpoint());
    }

    @Step("Connect Playwright browser via hub WS ({wsEndpoint})")
    public static Browser connect(Playwright playwright, String wsEndpoint) {
        if (wsEndpoint.contains("firefox")) {
            return playwright.firefox().connect(wsEndpoint);
        }
        if (wsEndpoint.contains("webkit")) {
            return playwright.webkit().connect(wsEndpoint);
        }
        return playwright.chromium().connect(wsEndpoint);
    }

    @Step("Close Playwright remote browser session")
    public static void close(Browser browser) {
        browser.close();
    }
}
