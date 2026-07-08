package api.hub;

import config.ConfigReader;

public final class HubRequest {

    private HubRequest() {
    }

    public static String baseUri() {
        var base = ConfigReader.resolveApiBaseUrl();
        return base.endsWith("/") ? base.substring(0, base.length() - 1) : base;
    }
}
