package resilience_test

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/mxschmitt/playwright-go"
	"github.com/stretchr/testify/require"
)

const connectedTimeout = 20 * time.Second

func TestMain(m *testing.M) {
	if err := playwright.Install(); err != nil {
		fmt.Fprintf(os.Stderr, "playwright install: %v\n", err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func openDashboard(t *testing.T, page playwright.Page, baseURL string) {
	t.Helper()
	_, err := page.Goto(baseURL+"/", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
	})
	require.NoError(t, err)
	waitConnected(t, page)
}

func waitConnected(t *testing.T, page playwright.Page) {
	t.Helper()
	waitStatusTileConnected(t, page, "#sse-status")
	waitStatusTileConnected(t, page, "#selenoid-status")
}

func waitStatusTileConnected(t *testing.T, page playwright.Page, selector string) {
	t.Helper()
	loc := page.Locator(selector)
	require.NoError(t, loc.WaitFor(playwright.LocatorWaitForOptions{
		State:   playwright.WaitForSelectorStateVisible,
		Timeout: playwright.Float(float64(connectedTimeout.Milliseconds())),
	}))
	deadline := time.Now().Add(connectedTimeout)
	for time.Now().Before(deadline) {
		if statusTileLooksConnected(loc) {
			return
		}
		time.Sleep(250 * time.Millisecond)
	}
	class, _ := loc.GetAttribute("class")
	text, _ := loc.InnerText()
	t.Fatalf("%s: expected connected tile within %s, class=%q text=%q", selector, connectedTimeout, class, text)
}

func statusTileLooksConnected(loc playwright.Locator) bool {
	class, err := loc.GetAttribute("class")
	if err != nil {
		return false
	}
	if strings.Contains(class, "status-tile--connected") {
		return true
	}
	text, err := loc.InnerText()
	if err != nil {
		return false
	}
	normalized := strings.ToUpper(strings.TrimSpace(text))
	return strings.Contains(normalized, "CONNECTED")
}

func stayConnected(t *testing.T, page playwright.Page, stable time.Duration) {
	t.Helper()
	step := 500 * time.Millisecond
	steps := int(stable / step)
	for i := 0; i < steps; i++ {
		waitConnected(t, page)
		time.Sleep(step)
	}
}
