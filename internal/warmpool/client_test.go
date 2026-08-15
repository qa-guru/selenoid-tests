package warmpool

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestClientHealthSlotsReserveRelease(t *testing.T) {
	var released string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/", "/health":
			_, _ = w.Write([]byte(`{"ok":true,"slots":2}`))
		case "/pool/slots":
			_, _ = w.Write([]byte(`[{"id":"pool-chrome-1","protocol":"webdriver","browser":"chrome","webdriverUrl":"http://127.0.0.1:14441/","reservedBy":null}]`))
		case "/pool/reserve":
			var in map[string]any
			_ = json.NewDecoder(r.Body).Decode(&in)
			if in["loopback"] != true {
				t.Errorf("loopback=%v", in["loopback"])
			}
			_, _ = w.Write([]byte(`{"ok":true,"slot":{"id":"pool-chrome-1","webdriverUrl":"http://127.0.0.1:14441/","webdriverUrlLoopback":"http://127.0.0.1:14441/","reservedBy":"hub-1"}}`))
		case "/pool/release":
			var in map[string]string
			_ = json.NewDecoder(r.Body).Decode(&in)
			released = in["slotId"]
			_, _ = w.Write([]byte(`{"ok":true,"slotId":"pool-chrome-1"}`))
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)

	c := New(srv.URL)
	if err := c.Ping(); err != nil {
		t.Fatalf("ping: %v", err)
	}
	h, err := c.Health()
	if err != nil || !h.OK || h.Slots != 2 {
		t.Fatalf("health: %+v %v", h, err)
	}
	root, err := c.Root()
	if err != nil || !root.OK {
		t.Fatalf("root: %+v %v", root, err)
	}
	slots, err := c.Slots()
	if err != nil || len(slots) != 1 || !slots[0].IsLoopback() {
		t.Fatalf("slots: %+v %v", slots, err)
	}
	slot, status, err := c.Reserve("webdriver", "chrome", "hub-1", true)
	if err != nil || status != 200 || slot.ID != "pool-chrome-1" {
		t.Fatalf("reserve: %d %+v %v", status, slot, err)
	}
	if err := c.Release(slot.ID); err != nil {
		t.Fatal(err)
	}
	if released != "pool-chrome-1" {
		t.Fatalf("released=%q", released)
	}
}

func TestClientReserve409(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusConflict)
		_, _ = w.Write([]byte(`{"error":"no available slots"}`))
	}))
	t.Cleanup(srv.Close)

	_, status, err := New(srv.URL).Reserve("webdriver", "chrome", "hub-1", true)
	if err != nil || status != 409 {
		t.Fatalf("status=%d err=%v", status, err)
	}
}

func TestParseError(t *testing.T) {
	if ParseError([]byte(`{"error":"slotId is required"}`)) != "slotId is required" {
		t.Fatal("ParseError")
	}
}

func TestClientPostAndGetStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/pool/reserve" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path == "/pool/release" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error":"slotId is required"}`))
			return
		}
		http.NotFound(w, r)
	}))
	t.Cleanup(srv.Close)

	c := New(srv.URL)
	status, body, err := c.Post("/pool/release", map[string]any{})
	if err != nil || status != 400 || ParseError(body) != "slotId is required" {
		t.Fatalf("post: %d %s %v", status, body, err)
	}
	status, _, err = c.GetStatus("/pool/reserve")
	if err != nil || status != http.StatusMethodNotAllowed {
		t.Fatalf("get reserve: %d %v", status, err)
	}
}

func TestBaseURLAndDefault(t *testing.T) {
	t.Setenv("WARM_POOL_URL", "")
	t.Setenv("SELENOID_WARM_POOL_URL", "")
	if BaseURL() != defaultBaseURL {
		t.Fatalf("default %q", BaseURL())
	}
	t.Setenv("SELENOID_WARM_POOL_URL", "http://alt:9090/")
	if BaseURL() != "http://alt:9090" {
		t.Fatalf("alt env %q", BaseURL())
	}
	t.Setenv("WARM_POOL_URL", " http://primary:9090/ ")
	if BaseURL() != "http://primary:9090" {
		t.Fatalf("primary env %q", BaseURL())
	}
	c := Default()
	if c == nil || c.http == nil {
		t.Fatal("Default")
	}
}

func TestClientErrorPaths(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/health", "/":
			http.Error(w, "down", http.StatusBadGateway)
		case "/pool/slots":
			http.Error(w, "nope", http.StatusInternalServerError)
		case "/pool/reserve":
			if r.Header.Get("X-Bad") == "json" {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{`))
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(strings.Repeat("e", 240)))
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)
	c := New(srv.URL)

	if _, err := c.Health(); err == nil {
		t.Fatal("health")
	}
	if _, err := c.Root(); err == nil {
		t.Fatal("root")
	}
	if _, err := c.Slots(); err == nil {
		t.Fatal("slots")
	}
	_, status, err := c.Reserve("webdriver", "chrome", "hub-1", true)
	if err == nil || status != http.StatusInternalServerError {
		t.Fatalf("reserve 500: %d %v", status, err)
	}
	if !strings.Contains(err.Error(), "…") {
		t.Fatalf("truncate: %v", err)
	}

	reqSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{`))
	}))
	t.Cleanup(reqSrv.Close)
	_, _, err = New(reqSrv.URL).Reserve("webdriver", "chrome", "hub-1", true)
	if err == nil {
		t.Fatal("reserve bad json")
	}

	dead := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	dead.Close()
	d := New(dead.URL)
	if _, _, err := d.Reserve("webdriver", "chrome", "hub-1", true); err == nil {
		t.Fatal("reserve down")
	}
	if _, _, err := d.Post("/x", map[string]any{}); err == nil {
		t.Fatal("post down")
	}
	if _, _, err := d.GetStatus("/x"); err == nil {
		t.Fatal("get down")
	}
}

func TestSlotIsLoopback(t *testing.T) {
	loop := Slot{WebdriverURL: "http://127.0.0.1:14441/"}
	if !loop.IsLoopback() {
		t.Fatal("expected loopback")
	}
	dns := Slot{WebdriverURL: "http://warm-chrome-1:4444/"}
	if dns.IsLoopback() {
		t.Fatal("docker-DNS must not be loopback")
	}
	pref := Slot{WebdriverURL: "http://warm-chrome-1:4444/", WebdriverURLLoopback: "http://127.0.0.1:14441/"}
	if !pref.IsLoopback() || pref.DialURL() != "http://127.0.0.1:14441/" {
		t.Fatalf("loopback field should win: %q", pref.DialURL())
	}
	empty := Slot{}
	if empty.IsLoopback() {
		t.Fatal("empty URL is not loopback")
	}
}

func TestTruncate(t *testing.T) {
	if truncate("ab", 5) != "ab" {
		t.Fatal("short")
	}
	if truncate("abcdef", 3) != "abc…" {
		t.Fatal("long")
	}
}
