package helpers;

import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.LinkedHashSet;
import java.util.List;
import java.util.Locale;
import java.util.Map;
import java.util.Set;
import org.openqa.selenium.json.Json;

/** Lightweight HAR 1.2 metrics for completeness comparisons. */
public final class HarStats {

    private static final Json JSON = new Json();

    public final String label;
    public final int entries;
    public final int httpEntries;
    public final int withStatus;
    public final int withResponseHeaders;
    public final int withRequestHeaders;
    public final int withContentSize;
    public final int withContentText;
    public final int withPostData;
    public final int withPositiveTime;
    public final Set<String> urls;
    public final String creator;
    public final Path path;

    private HarStats(
            String label,
            Path path,
            int entries,
            int httpEntries,
            int withStatus,
            int withResponseHeaders,
            int withRequestHeaders,
            int withContentSize,
            int withContentText,
            int withPostData,
            int withPositiveTime,
            Set<String> urls,
            String creator) {
        this.label = label;
        this.path = path;
        this.entries = entries;
        this.httpEntries = httpEntries;
        this.withStatus = withStatus;
        this.withResponseHeaders = withResponseHeaders;
        this.withRequestHeaders = withRequestHeaders;
        this.withContentSize = withContentSize;
        this.withContentText = withContentText;
        this.withPostData = withPostData;
        this.withPositiveTime = withPositiveTime;
        this.urls = urls;
        this.creator = creator;
    }

    @SuppressWarnings("unchecked")
    public static HarStats fromBytes(String label, Path path, byte[] harBytes) {
        Map<String, Object> root = JSON.toType(new String(harBytes, StandardCharsets.UTF_8), Map.class);
        Map<String, Object> log = (Map<String, Object>) root.getOrDefault("log", Map.of());
        List<Map<String, Object>> entries = (List<Map<String, Object>>) log.getOrDefault("entries", List.of());
        Map<String, Object> creatorMap = (Map<String, Object>) log.getOrDefault("creator", Map.of());
        String creator = String.valueOf(creatorMap.getOrDefault("name", ""));

        int withStatus = 0;
        int withResponseHeaders = 0;
        int withRequestHeaders = 0;
        int withContentSize = 0;
        int withContentText = 0;
        int withPostData = 0;
        int withPositiveTime = 0;
        int httpEntries = 0;
        Set<String> urls = new LinkedHashSet<>();

        for (Map<String, Object> entry : entries) {
            Map<String, Object> req = (Map<String, Object>) entry.getOrDefault("request", Map.of());
            Map<String, Object> resp = (Map<String, Object>) entry.getOrDefault("response", Map.of());
            Map<String, Object> content = (Map<String, Object>) resp.getOrDefault("content", Map.of());
            String url = String.valueOf(req.getOrDefault("url", ""));
            if (isHttpUrl(url)) {
                httpEntries++;
                urls.add(stripQuery(url));
            }
            if (intVal(resp.get("status")) > 0) {
                withStatus++;
            }
            if (listSize(resp.get("headers")) > 0) {
                withResponseHeaders++;
            }
            if (listSize(req.get("headers")) > 0) {
                withRequestHeaders++;
            }
            if (longVal(content.get("size")) > 0) {
                withContentSize++;
            }
            Object text = content.get("text");
            if (text instanceof String s && !s.isBlank()) {
                withContentText++;
            }
            if (req.get("postData") instanceof Map<?, ?> pd && !pd.isEmpty()) {
                withPostData++;
            }
            if (doubleVal(entry.get("time")) > 0) {
                withPositiveTime++;
            }
        }

        return new HarStats(
                label,
                path,
                entries.size(),
                httpEntries,
                withStatus,
                withResponseHeaders,
                withRequestHeaders,
                withContentSize,
                withContentText,
                withPostData,
                withPositiveTime,
                urls,
                creator);
    }

    public static HarStats fromFile(String label, Path path) throws Exception {
        return fromBytes(label, path, Files.readAllBytes(path));
    }

    public Map<String, Object> toRow() {
        Map<String, Object> row = new LinkedHashMap<>();
        row.put("label", label);
        row.put("creator", creator);
        row.put("entries", entries);
        row.put("httpEntries", httpEntries);
        row.put("uniqueHttpUrls", urls.size());
        row.put("withStatus", withStatus);
        row.put("withRequestHeaders", withRequestHeaders);
        row.put("withResponseHeaders", withResponseHeaders);
        row.put("withContentSize", withContentSize);
        row.put("withContentText", withContentText);
        row.put("withPostData", withPostData);
        row.put("withPositiveTime", withPositiveTime);
        row.put("path", path == null ? "" : path.toString());
        return row;
    }

    /** Fraction of baseline HTTP URLs also present in this HAR (0..1). */
    public double urlCoverageOf(HarStats baseline) {
        if (baseline.urls.isEmpty()) {
            return 1.0;
        }
        long hit = baseline.urls.stream().filter(urls::contains).count();
        return hit / (double) baseline.urls.size();
    }

    public static String formatTable(List<HarStats> stats) {
        List<String> lines = new ArrayList<>();
        lines.add(String.format(
                Locale.ROOT,
                "%-32s %8s %8s %8s %8s %8s %8s %8s %8s %8s",
                "label",
                "entries",
                "http",
                "urls",
                "status",
                "reqHdr",
                "resHdr",
                "size>0",
                "text",
                "time>0"));
        for (HarStats s : stats) {
            lines.add(String.format(
                    Locale.ROOT,
                    "%-32s %8d %8d %8d %8d %8d %8d %8d %8d %8d",
                    s.label,
                    s.entries,
                    s.httpEntries,
                    s.urls.size(),
                    s.withStatus,
                    s.withRequestHeaders,
                    s.withResponseHeaders,
                    s.withContentSize,
                    s.withContentText,
                    s.withPositiveTime));
        }
        return String.join("\n", lines);
    }

    private static boolean isHttpUrl(String url) {
        String u = url.toLowerCase(Locale.ROOT);
        return u.startsWith("http://") || u.startsWith("https://");
    }

    private static String stripQuery(String url) {
        String normalized = url
                .replace("host.docker.internal", "127.0.0.1")
                .replace("localhost", "127.0.0.1");
        int q = normalized.indexOf('?');
        return q < 0 ? normalized : normalized.substring(0, q);
    }

    private static int listSize(Object o) {
        return o instanceof List<?> list ? list.size() : 0;
    }

    private static int intVal(Object o) {
        if (o instanceof Number n) {
            return n.intValue();
        }
        return 0;
    }

    private static long longVal(Object o) {
        if (o instanceof Number n) {
            return n.longValue();
        }
        return 0L;
    }

    private static double doubleVal(Object o) {
        if (o instanceof Number n) {
            return n.doubleValue();
        }
        return 0.0;
    }
}
