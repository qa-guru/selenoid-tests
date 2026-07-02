package api.hub;

public record HubPingResponse(String uptime, String lastReloadTime, long numRequests, String version) {
}
