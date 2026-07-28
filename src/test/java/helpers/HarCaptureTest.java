package helpers;

import java.nio.charset.StandardCharsets;
import java.nio.file.Path;
import java.util.List;
import java.util.Map;
import java.util.logging.Level;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;
import org.openqa.selenium.logging.LogEntries;
import org.openqa.selenium.logging.LogEntry;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertTrue;

@Tag("unit")
class HarCaptureTest {

    @Test
    void toHarBuildsEntriesFromPerformanceLogs() {
        LogEntries entries = fixtureEntries();

        String har = HarCapture.toHar(entries);
        assertTrue(har.contains("1.2"), () -> "HAR missing version: " + har);
        assertTrue(har.contains("example.com"), () -> "HAR missing url: " + har);
        assertTrue(har.contains("200"), () -> "HAR missing status: " + har);
        assertTrue(har.contains("1280"), () -> "HAR missing loadingFinished size: " + har);
        assertTrue(HarCapture.supportsBrowser("chrome"));
        assertTrue(!HarCapture.supportsBrowser("firefox"));
    }

    @Test
    void toHarDefaultMetaOmitsContentText() {
        LogEntries entries = fixtureEntries();

        String harDefault = HarCapture.toHar(entries);
        String harMeta = HarCapture.toHar(entries, HarCapture.HarContentMode.META);

        HarStats statsDefault = HarStats.fromBytes(
                "unit-meta-default",
                Path.of("unit-meta-default.har"),
                harDefault.getBytes(StandardCharsets.UTF_8));
        HarStats statsMeta = HarStats.fromBytes(
                "unit-meta",
                Path.of("unit-meta.har"),
                harMeta.getBytes(StandardCharsets.UTF_8));

        assertEquals(0, statsDefault.withContentText, "default toHar must stay meta (no content.text)");
        assertEquals(0, statsMeta.withContentText, "explicit META must omit content.text");
        assertFalse(harDefault.contains("\"text\""), () -> "default HAR leaked content.text: " + harDefault);
        assertFalse(harMeta.contains("\"text\""), () -> "META HAR leaked content.text: " + harMeta);
    }

    @Test
    void toHarBodiesIncludesSyntheticContentText() {
        LogEntries entries = fixtureEntries();
        // Performance logs have no bodies; BODIES mode consumes CDP (or synthetic) payloads by requestId.
        Map<String, HarCapture.CapturedBody> bodies = Map.of(
                "r1", new HarCapture.CapturedBody("<html>synthetic-body</html>", false));

        String har = HarCapture.toHar(entries, HarCapture.HarContentMode.BODIES, bodies);
        HarStats stats = HarStats.fromBytes(
                "unit-bodies",
                Path.of("unit-bodies.har"),
                har.getBytes(StandardCharsets.UTF_8));

        assertTrue(stats.withContentText > 0, "BODIES + synthetic fixture must yield content.text");
        assertTrue(har.contains("synthetic-body"), () -> "HAR missing body text: " + har);
        assertFalse(har.contains("\"encoding\""), "plain text body must not set encoding=base64");
    }

    @Test
    void toHarBodiesBase64SetsEncoding() {
        LogEntries entries = fixtureEntries();
        Map<String, HarCapture.CapturedBody> bodies = Map.of(
                "r1", new HarCapture.CapturedBody("aGVsbG8=", true));

        String har = HarCapture.toHar(entries, HarCapture.HarContentMode.BODIES, bodies);

        assertTrue(har.contains("aGVsbG8="), () -> "HAR missing base64 body: " + har);
        assertTrue(har.contains("base64"), () -> "HAR missing encoding=base64: " + har);
    }

    @Test
    void toHarBodiesWithoutPayloadStaysMetaLike() {
        LogEntries entries = fixtureEntries();

        String har = HarCapture.toHar(entries, HarCapture.HarContentMode.BODIES, Map.of());
        HarStats stats = HarStats.fromBytes(
                "unit-bodies-empty",
                Path.of("unit-bodies-empty.har"),
                har.getBytes(StandardCharsets.UTF_8));

        assertEquals(0, stats.withContentText, "BODIES without CDP/synthetic payload must stay meta-like");
    }

    @Test
    void finishedRequestIdsFromLoadingFinished() {
        assertEquals(List.of("r1"), HarCapture.finishedRequestIds(fixtureEntries()));
    }

    private static LogEntries fixtureEntries() {
        long now = System.currentTimeMillis();
        String requestMsg = """
                {"message":{"method":"Network.requestWillBeSent","params":{"requestId":"r1","timestamp":1.0,"wallTime":1700000000.0,"request":{"url":"https://example.com/","method":"GET","headers":{"Accept":"*/*"}}}}}
                """.trim();
        String responseMsg = """
                {"message":{"method":"Network.responseReceived","params":{"requestId":"r1","response":{"status":200,"statusText":"OK","mimeType":"text/html","headers":{"content-type":"text/html"},"protocol":"http/1.1","encodedDataLength":42}}}}
                """.trim();
        String finishedMsg = """
                {"message":{"method":"Network.loadingFinished","params":{"requestId":"r1","timestamp":1.05,"encodedDataLength":1280}}}
                """.trim();

        return new LogEntries(List.of(
                new LogEntry(Level.INFO, now, requestMsg),
                new LogEntry(Level.INFO, now + 1, responseMsg),
                new LogEntry(Level.INFO, now + 2, finishedMsg)
        ));
    }
}
