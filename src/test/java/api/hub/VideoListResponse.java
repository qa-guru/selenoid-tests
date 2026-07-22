package api.hub;

import java.util.List;

public record VideoListResponse(List<String> videos, int total, int limit, int offset) {
}
