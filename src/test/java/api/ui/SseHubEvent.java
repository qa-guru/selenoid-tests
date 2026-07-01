package api.ui;

import api.hub.HubStatus;

import java.util.List;

public record SseHubEvent(HubStatus state, List<String> errors) {

    public boolean hasState() {
        return state != null;
    }

    public boolean hasErrors() {
        return errors != null && !errors.isEmpty();
    }
}
