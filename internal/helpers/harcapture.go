// Package helpers ports pure offline unit helpers from Java helpers.* (−live hub).
package helpers

import (
	"encoding/json"
	"strings"
	"time"
)

// HarContentMode controls HAR content.text depth (HarCapture.HarContentMode).
type HarContentMode string

const (
	HarMeta   HarContentMode = "META"
	HarBodies HarContentMode = "BODIES"
)

// CapturedBody is a CDP/synthetic response body keyed by requestId.
type CapturedBody struct {
	Text          string
	Base64Encoded bool
}

// SupportsBrowser reports whether client HAR capture works for the browser
// (Chrome/Edge/Chromium only).
func SupportsBrowser(browser string) bool {
	if strings.TrimSpace(browser) == "" {
		return false
	}
	b := strings.ToLower(browser)
	return strings.Contains(b, "chrome") || strings.Contains(b, "edge") || b == "chromium"
}

// FinishedRequestIds returns Network.loadingFinished requestIds in order.
func FinishedRequestIds(logMessages []string) []string {
	var ids []string
	for _, msg := range logMessages {
		method, params, ok := parseCDPMessage(msg)
		if !ok || method != "Network.loadingFinished" {
			continue
		}
		id := stringVal(params["requestId"])
		if id != "" {
			ids = append(ids, id)
		}
	}
	return ids
}

// ToHar builds HAR 1.2 JSON from Performance-log style CDP message strings.
func ToHar(logMessages []string, mode HarContentMode, bodies map[string]CapturedBody) string {
	if mode == "" {
		mode = HarMeta
	}
	if bodies == nil {
		bodies = map[string]CapturedBody{}
	}

	requests := map[string]map[string]any{}
	responses := map[string]map[string]any{}
	finishedMs := map[string]float64{}
	encodedBytes := map[string]int64{}
	var order []string
	wallStart := nan

	for _, msg := range logMessages {
		method, params, ok := parseCDPMessage(msg)
		if !ok {
			continue
		}
		switch method {
		case "Network.requestWillBeSent":
			id := stringVal(params["requestId"])
			reqRaw, ok := params["request"].(map[string]any)
			if id == "" || !ok {
				continue
			}
			ts := doubleVal(params["timestamp"])
			if wallStart == nan && ts != nan {
				wallStart = ts
			}
			stored := map[string]any{
				"url":       stringVal(reqRaw["url"]),
				"method":    stringVal(reqRaw["method"]),
				"headers":   headersOrEmpty(reqRaw["headers"]),
				"timestamp": ts,
			}
			if _, has := params["wallTime"]; has {
				stored["wallTime"] = doubleVal(params["wallTime"])
			}
			if _, exists := requests[id]; !exists {
				order = append(order, id)
			}
			requests[id] = stored
		case "Network.responseReceived":
			id := stringVal(params["requestId"])
			respRaw, ok := params["response"].(map[string]any)
			if id == "" || !ok {
				continue
			}
			responses[id] = respRaw
		case "Network.loadingFinished":
			id := stringVal(params["requestId"])
			if id == "" {
				continue
			}
			finishedMs[id] = doubleVal(params["timestamp"])
			if _, has := params["encodedDataLength"]; has {
				encodedBytes[id] = longVal(params["encodedDataLength"])
			}
		}
	}

	var harEntries []map[string]any
	for _, id := range order {
		req := requests[id]
		if req == nil {
			continue
		}
		resp := responses[id]
		if resp == nil {
			resp = map[string]any{}
		}
		start := doubleVal(req["timestamp"])
		end, hasEnd := finishedMs[id]
		if !hasEnd {
			end = start
		}
		timeMs := 0.0
		if start != nan && end != nan && end >= start {
			timeMs = (end - start) * 1000.0
		}

		startedDateTimeMs := time.Now().UnixMilli()
		if _, has := req["wallTime"]; has {
			startedDateTimeMs = int64(doubleVal(req["wallTime"]) * 1000.0)
		} else if wallStart != nan && end != nan {
			startedDateTimeMs = time.Now().UnixMilli() - int64((end-wallStart)*1000.0)
		}

		var body *CapturedBody
		if mode == HarBodies {
			if b, ok := bodies[id]; ok {
				body = &b
			}
		}

		entry := map[string]any{
			"startedDateTime": time.UnixMilli(startedDateTimeMs).UTC().Format(time.RFC3339Nano),
			"time":            timeMs,
			"request":         harRequest(req),
			"response":        harResponse(resp, encodedBytes[id], mode, body),
			"cache":           map[string]any{},
			"timings":         timings(timeMs),
		}
		harEntries = append(harEntries, entry)
	}

	log := map[string]any{
		"version": "1.2",
		"creator": map[string]any{
			"name":    "zero-design-system HarCapture",
			"version": "1",
		},
		"pages": []map[string]any{{
			"startedDateTime": time.Now().UTC().Format(time.RFC3339Nano),
			"id":              "page_1",
			"title":           "selenide-har",
			"pageTimings": map[string]any{
				"onContentLoad": -1,
				"onLoad":        -1,
			},
		}},
		"entries": harEntries,
	}
	raw, _ := json.Marshal(map[string]any{"log": log})
	return string(raw)
}

const nan = -1e300 // sentinel; CDP timestamps are never this low

func parseCDPMessage(raw string) (method string, params map[string]any, ok bool) {
	var root map[string]any
	if err := json.Unmarshal([]byte(raw), &root); err != nil {
		return "", nil, false
	}
	var message map[string]any
	switch m := root["message"].(type) {
	case string:
		if err := json.Unmarshal([]byte(m), &message); err != nil {
			return "", nil, false
		}
	case map[string]any:
		message = m
	default:
		if _, hasMethod := root["method"]; hasMethod {
			message = root
		} else {
			return "", nil, false
		}
	}
	method = stringVal(message["method"])
	params, _ = message["params"].(map[string]any)
	if method == "" || params == nil {
		return "", nil, false
	}
	return method, params, true
}

func harRequest(req map[string]any) map[string]any {
	method := stringVal(req["method"])
	if method == "" {
		method = "GET"
	}
	return map[string]any{
		"method":      method,
		"url":         stringVal(req["url"]),
		"httpVersion": "HTTP/1.1",
		"cookies":     []any{},
		"headers":     headerList(req["headers"]),
		"queryString": []any{},
		"headersSize": -1,
		"bodySize":    -1,
	}
}

func harResponse(resp map[string]any, finishedEncoded int64, mode HarContentMode, body *CapturedBody) map[string]any {
	protocol := stringVal(resp["protocol"])
	if protocol == "" {
		protocol = "HTTP/1.1"
	}
	size := finishedEncoded
	if size == 0 {
		size = longVal(resp["encodedDataLength"])
	}
	content := map[string]any{
		"size":     size,
		"mimeType": stringVal(resp["mimeType"]),
	}
	if mode == HarBodies && body != nil && body.Text != "" {
		content["text"] = body.Text
		if body.Base64Encoded {
			content["encoding"] = "base64"
		}
	}
	return map[string]any{
		"status":      intVal(resp["status"]),
		"statusText":  stringVal(resp["statusText"]),
		"httpVersion": protocol,
		"cookies":     []any{},
		"headers":     headerList(resp["headers"]),
		"content":     content,
		"redirectURL": "",
		"headersSize": -1,
		"bodySize":    -1,
	}
}

func headerList(headersObj any) []map[string]string {
	headers, ok := headersObj.(map[string]any)
	if !ok {
		return []map[string]string{}
	}
	out := make([]map[string]string, 0, len(headers))
	for k, v := range headers {
		val := ""
		if v != nil {
			val = stringVal(v)
		}
		out = append(out, map[string]string{"name": k, "value": val})
	}
	return out
}

func timings(totalMs float64) map[string]any {
	wait := totalMs
	if wait < 0 {
		wait = 0
	}
	return map[string]any{
		"blocked": -1,
		"dns":     -1,
		"connect": -1,
		"ssl":     -1,
		"send":    0,
		"wait":    wait,
		"receive": 0,
	}
}

func headersOrEmpty(v any) any {
	if m, ok := v.(map[string]any); ok {
		return m
	}
	return map[string]any{}
}

func stringVal(o any) string {
	if o == nil {
		return ""
	}
	switch v := o.(type) {
	case string:
		return v
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return ""
		}
		s := string(b)
		if len(s) >= 2 && s[0] == '"' {
			return strings.Trim(s, `"`)
		}
		return s
	}
}

func doubleVal(o any) float64 {
	switch v := o.(type) {
	case float64:
		return v
	case json.Number:
		f, _ := v.Float64()
		return f
	case int:
		return float64(v)
	case int64:
		return float64(v)
	case nil:
		return nan
	default:
		return nan
	}
}

func intVal(o any) int {
	switch v := o.(type) {
	case float64:
		return int(v)
	case int:
		return v
	case int64:
		return int(v)
	default:
		return 0
	}
}

func longVal(o any) int64 {
	switch v := o.(type) {
	case float64:
		return int64(v)
	case int:
		return int64(v)
	case int64:
		return v
	default:
		return 0
	}
}
