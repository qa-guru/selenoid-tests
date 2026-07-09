package api.ui;

import config.ConfigReader;

public final class UiRequest {

    private UiRequest() {
    }

    public static String baseUri() {
        var base = ConfigReader.resolveUiUrl();
        return base.endsWith("/") ? base.substring(0, base.length() - 1) : base;
    }
}
