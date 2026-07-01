package api.hub;

import java.util.Map;

public record HubStatus(
        int total,
        int used,
        int queued,
        int pending,
        Map<String, Object> browsers
) {
}
