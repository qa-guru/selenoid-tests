package allure;

import com.codeborne.selenide.Selenide;
import io.qameta.allure.Allure;
import io.qameta.allure.Attachment;
import java.io.ByteArrayInputStream;
import java.nio.file.Files;
import java.nio.file.Path;
import org.openqa.selenium.OutputType;

public final class Attachments {

    private Attachments() {
    }

    @Attachment(value = "{attachName}", type = "image/png")
    public static byte[] screenshot(String attachName) {
        return Selenide.screenshot(OutputType.BYTES);
    }

    @Attachment(value = "{attachName}", type = "image/png")
    public static byte[] png(String attachName, byte[] bytes) {
        return bytes;
    }

    /** Attach a HAR file produced by Playwright {@code recordHar} (or any client capture). */
    public static void har(Path harFile) {
        if (harFile == null || !Files.isRegularFile(harFile)) {
            return;
        }
        try {
            Allure.addAttachment(
                    "HAR",
                    "application/json",
                    new ByteArrayInputStream(Files.readAllBytes(harFile)),
                    ".har");
        } catch (Exception ignored) {
            // never fail the test on attachment I/O
        }
    }
}
