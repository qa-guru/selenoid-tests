package warmpool

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
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
}
