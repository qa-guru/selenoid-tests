// Package stack controls the local Selenoid hub/UI process (Java StackHelper parity).
package stack

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/qa-guru/selenoid-tests/internal/config"
)

var httpClient = &http.Client{Timeout: 3 * time.Second}

// EnsureRunning starts hub and UI when their health endpoints are down (default :4444 / :8080).
func EnsureRunning() error {
	return EnsureControllable(nil)
}

// EnsureControllable (re)starts the ci-bin hub/UI pair on profile ports so kill/restart tests
// control the same backend the UI proxies (avoids mismatch with persistent stands on other ports).
func EnsureControllable(cfg *config.Config) error {
	hubStatus, uiStatus := defaultStatusURLs()
	if cfg != nil {
		if u, err := statusURL(cfg.HubURL); err == nil {
			hubStatus = u
		}
		if u, err := statusURL(cfg.UIURL); err == nil {
			uiStatus = u
		}
	}

	hubPort, err := portFromStatusURL(hubStatus)
	if err != nil {
		return err
	}
	uiPort, err := portFromStatusURL(uiStatus)
	if err != nil {
		return err
	}

	killPort(uiPort)
	killPort(hubPort)
	if err := waitEndpointDown(uiStatus, 10*time.Second); err != nil {
		return err
	}
	if err := waitEndpointDown(hubStatus, 10*time.Second); err != nil {
		return err
	}

	if err := StartHubDetached(); err != nil {
		return err
	}
	if err := waitEndpointUp(hubStatus, 30*time.Second); err != nil {
		return err
	}
	if err := StartUIDetached(); err != nil {
		return err
	}
	return waitEndpointUp(uiStatus, 30*time.Second)
}

// KillHub stops the hub listener for cfg (or :4444 when cfg is nil).
func KillHub(cfg *config.Config) error {
	return killConfiguredEndpoint(cfg, true)
}

// KillUI stops the UI listener for cfg (or :8080 when cfg is nil).
func KillUI(cfg *config.Config) error {
	return killConfiguredEndpoint(cfg, false)
}

func killConfiguredEndpoint(cfg *config.Config, hub bool) error {
	raw := defaultHubStatusURL
	if hub {
		if cfg != nil {
			if u, err := statusURL(cfg.HubURL); err == nil {
				raw = u
			}
		}
	} else {
		raw = defaultUIStatusURL
		if cfg != nil {
			if u, err := statusURL(cfg.UIURL); err == nil {
				raw = u
			}
		}
	}
	port, err := portFromStatusURL(raw)
	if err != nil {
		return err
	}
	return killPort(port)
}

// StartHubDetached launches the CI or dev hub binary in the background.
func StartHubDetached() error {
	root, err := moduleRoot()
	if err != nil {
		return err
	}
	ciBin := resolveCIHubBinary(root)
	if ciBin != "" {
		browsers := filepath.Join(root, "fixtures", "ci-browsers.json")
		return execDetached(root, []string{
			ciBin,
			"-conf", browsers,
			"-limit", "3",
			"-video-recorder-image", "qaguru/video-recorder:latest",
		})
	}
	devStart := filepath.Join(root, "..", "dev", "scripts", "start-selenoid.sh")
	if fileExists(devStart) {
		return execDetached(root, []string{"bash", devStart})
	}
	return fmt.Errorf("no hub binary: missing %s and %s", filepath.Join(root, "build", "ci-bin", "selenoid"), devStart)
}

// StartUIDetached launches the CI or dev UI binary in the background.
func StartUIDetached() error {
	root, err := moduleRoot()
	if err != nil {
		return err
	}
	ciUI := filepath.Join(root, "build", "ci-bin", "selenoid-ui")
	browsers := filepath.Join(root, "fixtures", "ci-browsers.json")
	if isExecutable(ciUI) && fileExists(browsers) {
		return execDetached(root, []string{
			ciUI,
			"-listen", ":8080",
			"-selenoid-uri", "http://127.0.0.1:4444",
			"-browsers-conf", browsers,
			"-period", "4s",
		})
	}
	devStart := filepath.Join(root, "..", "dev", "scripts", "start-selenoid-ui.sh")
	if fileExists(devStart) {
		return execDetached(root, []string{"bash", devStart})
	}
	return fmt.Errorf("no UI binary: missing %s and %s", ciUI, devStart)
}

// WaitForHubReady polls hub /status until HTTP 200 or timeout.
func WaitForHubReady(cfg *config.Config, timeout time.Duration) error {
	raw := defaultHubStatusURL
	if cfg != nil {
		if u, err := statusURL(cfg.HubURL); err == nil {
			raw = u
		}
	}
	return waitEndpointUp(raw, timeout)
}

// WaitForHubDown polls until hub /status stops responding.
func WaitForHubDown(cfg *config.Config, timeout time.Duration) error {
	raw := defaultHubStatusURL
	if cfg != nil {
		if u, err := statusURL(cfg.HubURL); err == nil {
			raw = u
		}
	}
	return waitEndpointDown(raw, timeout)
}

// WaitForUIReady polls UI /status until HTTP 200 or timeout.
func WaitForUIReady(cfg *config.Config, timeout time.Duration) error {
	raw := defaultUIStatusURL
	if cfg != nil {
		if u, err := statusURL(cfg.UIURL); err == nil {
			raw = u
		}
	}
	return waitEndpointUp(raw, timeout)
}

// HubResponds returns true when default hub /status returns HTTP 200.
func HubResponds() bool {
	return endpointResponds(defaultHubStatusURL)
}

// UIResponds returns true when default UI /status returns HTTP 200.
func UIResponds() bool {
	return endpointResponds(defaultUIStatusURL)
}

const (
	defaultHubStatusURL = "http://127.0.0.1:4444/status"
	defaultUIStatusURL  = "http://127.0.0.1:8080/status"
)

func defaultStatusURLs() (hub, ui string) {
	return defaultHubStatusURL, defaultUIStatusURL
}

func statusURL(base string) (string, error) {
	base = strings.TrimSpace(base)
	if base == "" {
		return "", fmt.Errorf("empty url")
	}
	u, err := url.Parse(base)
	if err != nil {
		return "", err
	}
	if u.Scheme == "" {
		u.Scheme = "http"
	}
	u.Path = "/status"
	u.RawQuery = ""
	u.Fragment = ""
	return u.String(), nil
}

func portFromStatusURL(raw string) (int, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return 0, err
	}
	if u.Port() != "" {
		var port int
		_, err := fmt.Sscanf(u.Port(), "%d", &port)
		if err != nil {
			return 0, err
		}
		return port, nil
	}
	switch u.Scheme {
	case "https":
		return 443, nil
	default:
		return 80, nil
	}
}

func killPort(port int) error {
	return execShell(fmt.Sprintf("lsof -nP -iTCP:%d -sTCP:LISTEN -t | xargs kill -9 2>/dev/null || true", port))
}

func waitEndpointUp(raw string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if endpointResponds(raw) {
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return fmt.Errorf("endpoint did not become ready: %s within %s", raw, timeout)
}

func waitEndpointDown(raw string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if !endpointResponds(raw) {
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return fmt.Errorf("endpoint still responding: %s after %s", raw, timeout)
}

func endpointResponds(raw string) bool {
	resp, err := httpClient.Get(raw)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

func resolveCIHubBinary(root string) string {
	if fromEnv := strings.TrimSpace(os.Getenv("SELENOID_BIN")); fromEnv != "" && isExecutable(fromEnv) {
		return fromEnv
	}
	ciBin := filepath.Join(root, "build", "ci-bin", "selenoid")
	browsers := filepath.Join(root, "fixtures", "ci-browsers.json")
	if isExecutable(ciBin) && fileExists(browsers) {
		return ciBin
	}
	return ""
}

func moduleRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	dir := wd
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found from %s", wd)
		}
		dir = parent
	}
}

func execShell(command string) error {
	cmd := exec.Command("bash", "-lc", command)
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}

func execDetached(dir string, command []string) error {
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Dir = dir
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Start()
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func isExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil || info.IsDir() {
		return false
	}
	return info.Mode()&0o111 != 0
}
