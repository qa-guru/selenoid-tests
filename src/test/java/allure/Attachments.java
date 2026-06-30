package allure;

import com.codeborne.selenide.Selenide;
import io.qameta.allure.Attachment;
import org.openqa.selenium.OutputType;

public final class Attachments {

    private Attachments() {
    }

    @Attachment(value = "{attachName}", type = "image/png")
    public static byte[] screenshot(String attachName) {
        return Selenide.screenshot(OutputType.BYTES);
    }
}
