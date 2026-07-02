package api.hub;

import java.util.List;

public record ChromeOptions(List<String> args) {

    static ChromeOptions dockerHeadless() {
        return new ChromeOptions(List.of(
                "--headless=new",
                "--no-sandbox",
                "--disable-dev-shm-usage"
        ));
    }
}
