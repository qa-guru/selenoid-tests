package api.hub;

import java.util.List;

public record FirefoxOptions(List<String> args) {

    public static FirefoxOptions dockerHeadless() {
        return new FirefoxOptions(List.of("-headless"));
    }
}
