package api.hub;

import java.util.List;

public record EdgeOptions(List<String> args) {

    public static EdgeOptions dockerHeadless() {
        return new EdgeOptions(List.of(
                "--headless=new",
                "--no-sandbox",
                "--disable-dev-shm-usage"
        ));
    }
}
