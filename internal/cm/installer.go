package cm

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/qa-guru/selenoid-tests/internal/config"
)

// Installer drives cm selenoid / selenoid-ui lifecycle (Java CmInstallerHelper).
type Installer struct {
	cfg           *config.Config
	configDir     string
	projectRoot   string
	workspaceRoot string
	cmBinary      string
	browsersJSON  string
}

// WithTempConfigDir creates an isolated CM config directory.
func WithTempConfigDir(cfg *config.Config) (*Installer, error) {
	root, err := moduleRoot()
	if err != nil {
		return nil, err
	}
	dir, err := os.MkdirTemp("", "cm-installer-")
	if err != nil {
		return nil, err
	}
	parent := filepath.Dir(root)
	return &Installer{
		cfg:           cfg,
		configDir:     dir,
		projectRoot:   root,
		workspaceRoot: parent,
	}, nil
}

// ConfigDir returns the temp CM configuration directory.
func (i *Installer) ConfigDir() string { return i.configDir }

// BrowsersJSONPath is configDir/browsers.json after configure.
func (i *Installer) BrowsersJSONPath() string {
	return filepath.Join(i.configDir, "browsers.json")
}

// DeleteConfigDir removes the temp config tree (best effort).
func (i *Installer) DeleteConfigDir() {
	_ = os.RemoveAll(i.configDir)
}

// StopAll stops CM-managed hub/UI and releases published ports.
func (i *Installer) StopAll() {
	i.stopHubQuietly()
	i.stopUiQuietly()
	forceRemoveCmContainers()
	releasePublishedPorts(i.cfg.CmHubPort, i.cfg.CmUiPort)
}

// Configure runs cm selenoid configure -n.
func (i *Installer) Configure() (RunResult, error) {
	if err := i.ensureLinuxBinary("selenoid", i.cfg.CmSelenoidBinary, i.selenoidRepoDir()); err != nil {
		return RunResult{}, err
	}
	if err := i.ensureLinuxBinary("selenoid-ui", i.cfg.CmSelenoidUiBinary, i.selenoidUiRepoDir()); err != nil {
		return RunResult{}, err
	}
	return i.runSelenoid("configure",
		"-c", i.configDir,
		"-p", fmt.Sprint(i.cfg.CmHubPort),
		"-n",
		"-j", i.resolvedBrowsersJSON(),
		"--selenoid-binary", filepath.Join(i.configDir, "bin", "selenoid"),
		"--selenoid-ui-binary", filepath.Join(i.configDir, "bin", "selenoid-ui"),
	)
}

// StartHub runs cm selenoid start -f.
func (i *Installer) StartHub() (RunResult, error) {
	if err := i.ensureLinuxBinary("selenoid", i.cfg.CmSelenoidBinary, i.selenoidRepoDir()); err != nil {
		return RunResult{}, err
	}
	return i.runSelenoid("start",
		"-f",
		"-c", i.configDir,
		"-p", fmt.Sprint(i.cfg.CmHubPort),
		"-n",
		"-j", i.resolvedBrowsersJSON(),
		"--selenoid-binary", filepath.Join(i.configDir, "bin", "selenoid"),
		"--selenoid-ui-binary", filepath.Join(i.configDir, "bin", "selenoid-ui"),
	)
}

// StartUi runs cm selenoid-ui start.
func (i *Installer) StartUi() (RunResult, error) {
	if err := i.ensureLinuxBinary("selenoid-ui", i.cfg.CmSelenoidUiBinary, i.selenoidUiRepoDir()); err != nil {
		return RunResult{}, err
	}
	return i.runSelenoidUi("start",
		"-c", i.configDir,
		"-p", fmt.Sprint(i.cfg.CmUiPort),
	)
}

// StatusHub runs cm selenoid status.
func (i *Installer) StatusHub() RunResult {
	r, _ := i.runSelenoid("status", "-c", i.configDir, "-p", fmt.Sprint(i.cfg.CmHubPort))
	return r
}

// StatusUi runs cm selenoid-ui status.
func (i *Installer) StatusUi() RunResult {
	r, _ := i.runSelenoidUi("status", "-c", i.configDir, "-p", fmt.Sprint(i.cfg.CmUiPort))
	return r
}

// WaitForHubReady polls GET {cmHubUrl}status until HTTP 200.
func (i *Installer) WaitForHubReady(timeout time.Duration) error {
	return i.waitForStatusReady(i.cfg.ResolveCmHubURL(), timeout, "hub", i.StatusHub, "selenoid")
}

// WaitForUiReady polls GET {cmUiUrl}status until HTTP 200.
func (i *Installer) WaitForUiReady(timeout time.Duration) error {
	return i.waitForStatusReady(i.cfg.ResolveCmUiURL(), timeout, "UI", i.StatusUi, "selenoid-ui")
}

func (i *Installer) cmBinaryPath() (string, error) {
	if i.cmBinary != "" {
		return i.cmBinary, nil
	}
	p, err := resolveExecutable(i.projectRoot, i.cfg.CmBinaryPath, "cm binary")
	if err != nil {
		return "", err
	}
	i.cmBinary = p
	return p, nil
}

func (i *Installer) resolvedBrowsersJSON() string {
	p, err := i.browsersJSONPath()
	if err != nil {
		panic(err)
	}
	return p
}

func (i *Installer) browsersJSONPath() (string, error) {
	if i.browsersJSON != "" {
		return i.browsersJSON, nil
	}
	p, err := resolveExisting(i.projectRoot, i.cfg.CmBrowsersJSON, "browsers.json")
	if err != nil {
		return "", err
	}
	i.browsersJSON = p
	return p, nil
}

func (i *Installer) selenoidRepoDir() string {
	if d := filepath.Join(i.projectRoot, "repos", "selenoid"); dirExists(d) {
		return d
	}
	return filepath.Join(i.workspaceRoot, "selenoid")
}

func (i *Installer) selenoidUiRepoDir() string {
	if d := filepath.Join(i.projectRoot, "repos", "selenoid-ui"); dirExists(d) {
		return d
	}
	return filepath.Join(i.workspaceRoot, "selenoid-ui")
}

func (i *Installer) ensureLinuxBinary(binaryName, localProperty, sourceRepo string) error {
	target := filepath.Join(i.configDir, "bin", binaryName)
	if isExecutable(target) {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return err
	}
	if i.cfg.CmUseLocalBinaries && isLinuxHost() {
		local, err := resolveExecutable(i.projectRoot, localProperty, binaryName+" binary")
		if err != nil {
			return err
		}
		if err := copyFile(local, target); err != nil {
			return err
		}
		return makeExecutable(target)
	}
	if binaryName == "selenoid-ui" {
		if err := runGoGenerate(sourceRepo); err != nil {
			return err
		}
	}
	return crossCompileLinuxBinary(sourceRepo, target)
}

func (i *Installer) runSelenoid(subcommand string, args ...string) (RunResult, error) {
	return i.run("selenoid", subcommand, args...)
}

func (i *Installer) runSelenoidUi(subcommand string, args ...string) (RunResult, error) {
	return i.run("selenoid-ui", subcommand, args...)
}

func (i *Installer) run(scope, subcommand string, args ...string) (RunResult, error) {
	binary, err := i.cmBinaryPath()
	if err != nil {
		return RunResult{}, err
	}
	cmdArgs := append([]string{scope, subcommand}, args...)
	return execCm(binary, cmdArgs)
}

func (i *Installer) stopHubQuietly() {
	_, _ = i.runSelenoid("stop", "-c", i.configDir, "-p", fmt.Sprint(i.cfg.CmHubPort))
}

func (i *Installer) stopUiQuietly() {
	_, _ = i.runSelenoidUi("stop", "-c", i.configDir, "-p", fmt.Sprint(i.cfg.CmUiPort))
}

func (i *Installer) waitForStatusReady(baseURL string, timeout time.Duration, label string, statusFn func() RunResult, containerName string) error {
	deadline := time.Now().Add(timeout)
	statusURI := strings.TrimRight(baseURL, "/") + "/status"
	client := &http.Client{Timeout: 3 * time.Second}
	for time.Now().Before(deadline) {
		if statusResponds(client, statusURI) {
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	lastStatus := statusFn().Output
	logs := dockerLogsTail(containerName)
	msg := fmt.Sprintf("%s did not become ready at %s within %s", label, statusURI, timeout)
	if lastStatus != "" {
		msg += "\ncm status:\n" + lastStatus
	}
	if logs != "" {
		msg += "\ndocker logs " + containerName + ":\n" + logs
	}
	return fmt.Errorf("%s", msg)
}

// Version runs cm version.
func Version(cfg *config.Config) (RunResult, error) {
	root, err := moduleRoot()
	if err != nil {
		return RunResult{}, err
	}
	binary, err := resolveExecutable(root, cfg.CmBinaryPath, "cm binary")
	if err != nil {
		return RunResult{}, err
	}
	return execCm(binary, []string{"version"})
}

// Help runs cm --help.
func Help(cfg *config.Config) (RunResult, error) {
	root, err := moduleRoot()
	if err != nil {
		return RunResult{}, err
	}
	binary, err := resolveExecutable(root, cfg.CmBinaryPath, "cm binary")
	if err != nil {
		return RunResult{}, err
	}
	return execCm(binary, []string{"--help"})
}

func execCm(binary string, args []string) (RunResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, binary, args...)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return RunResult{}, fmt.Errorf("cm %s timed out after 5 minutes", strings.Join(args, " "))
		}
		if exit, ok := err.(*exec.ExitError); ok {
			return RunResult{ExitCode: exit.ExitCode(), Output: buf.String()}, nil
		}
		return RunResult{}, fmt.Errorf("failed to run cm %s: %w", strings.Join(args, " "), err)
	}
	return RunResult{ExitCode: 0, Output: buf.String()}, nil
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

func resolveExecutable(projectRoot, path, label string) (string, error) {
	resolved := path
	if !filepath.IsAbs(resolved) {
		resolved = filepath.Join(projectRoot, resolved)
	}
	resolved = filepath.Clean(resolved)
	if !isExecutable(resolved) {
		return "", fmt.Errorf("%s not found or not executable: %s", label, resolved)
	}
	return resolved, nil
}

func resolveExisting(projectRoot, path, label string) (string, error) {
	resolved := path
	if !filepath.IsAbs(resolved) {
		resolved = filepath.Join(projectRoot, resolved)
	}
	resolved = filepath.Clean(resolved)
	info, err := os.Stat(resolved)
	if err != nil || info.IsDir() {
		return "", fmt.Errorf("%s not found: %s", label, resolved)
	}
	return resolved, nil
}

func isExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	if info.IsDir() {
		return false
	}
	return info.Mode()&0o111 != 0
}

func isLinuxHost() bool {
	return strings.Contains(strings.ToLower(runtime.GOOS), "linux")
}

func linuxGoArch() string {
	if runtime.GOOS == "darwin" && (runtime.GOARCH == "arm64" || strings.Contains(runtime.GOARCH, "arm")) {
		return "amd64"
	}
	if strings.Contains(runtime.GOARCH, "arm") {
		return "arm64"
	}
	return "amd64"
}

func makeExecutable(path string) error {
	return os.Chmod(path, 0o755)
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o755)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func runGoGenerate(sourceRepo string) error {
	if !dirExists(sourceRepo) {
		return fmt.Errorf("source repo not found: %s", sourceRepo)
	}
	goPath, err := exec.Command("go", "env", "GOPATH").Output()
	if err != nil {
		return fmt.Errorf("go env GOPATH: %w", err)
	}
	pathEnv := strings.TrimSpace(string(goPath)) + "/bin:" + os.Getenv("PATH")
	install := exec.Command("go", "install", "github.com/rakyll/statik@latest")
	install.Env = append(os.Environ(), "PATH="+pathEnv)
	_, _ = install.CombinedOutput()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	gen := exec.CommandContext(ctx, "go", "generate", ".")
	gen.Dir = sourceRepo
	gen.Env = append(os.Environ(), "PATH="+pathEnv)
	out, err := gen.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("go generate timed out in %s", sourceRepo)
	}
	if err != nil {
		return fmt.Errorf("go generate failed in %s: %s", sourceRepo, string(out))
	}
	return nil
}

func crossCompileLinuxBinary(sourceRepo, target string) error {
	if !dirExists(sourceRepo) {
		return fmt.Errorf("source repo not found: %s", sourceRepo)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, "go", "build", "-o", target, ".")
	cmd.Dir = sourceRepo
	cmd.Env = append(os.Environ(),
		"GOOS=linux",
		"GOARCH="+linuxGoArch(),
		"CGO_ENABLED=0",
	)
	out, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("go build timed out in %s", sourceRepo)
	}
	if err != nil {
		return fmt.Errorf("go build failed in %s: %s", sourceRepo, string(out))
	}
	return makeExecutable(target)
}

func forceRemoveCmContainers() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "docker", "rm", "-f", "selenoid", "selenoid-ui")
	_ = cmd.Run()
}

func releasePublishedPorts(ports ...int) {
	for _, port := range ports {
		stopDockerPublishPort(port)
	}
}

func stopDockerPublishPort(port int) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	list := exec.CommandContext(ctx, "docker", "ps", "-q", "--filter", fmt.Sprintf("publish=%d", port))
	out, err := list.Output()
	if err != nil || len(bytes.TrimSpace(out)) == 0 {
		return
	}
	for _, id := range strings.Fields(string(out)) {
		stop := exec.CommandContext(ctx, "docker", "stop", id)
		_ = stop.Run()
	}
}

func dockerLogsTail(containerName string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "docker", "logs", "--tail", "50", containerName)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func statusResponds(client *http.Client, statusURI string) bool {
	resp, err := client.Get(statusURI)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)
	return resp.StatusCode == http.StatusOK
}
