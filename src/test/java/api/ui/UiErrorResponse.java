package api.ui;

import java.util.List;

public record UiErrorResponse(List<UiErrorMessage> errors, String version) {
}
