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
        return playwright.chromium().connect(ConfigReader.resolvePlaywrightWsEndpoint());
    }

    @Step("Close Playwright remote browser session")
    public static void close(Browser browser) {
        browser.close();
    }
}
