package playwrightdebug

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

//go:embed bootstrap.cjs
var bootstrapScript string

// RunPlanResult describes a planned file-mode invocation without executing node.
type RunPlanResult struct {
	ScriptPath        string   `json:"script_path"`
	ScriptArgs        []string `json:"script_args"`
	CacheDir          string   `json:"cache_dir"`
	PlaywrightVersion string   `json:"playwright_version,omitempty"`
	NodeArgv          []string `json:"node_argv"`
	BootstrapWritten  bool     `json:"bootstrap_written"`
	Error             string   `json:"error,omitempty"`
}

// BuildRunPlan validates the script path and records the node argv that RunFile would use.
func BuildRunPlan(opts RunOptions) (RunPlanResult, error) {
	if err := ValidateScriptFile(opts.ScriptPath); err != nil {
		return RunPlanResult{Error: err.Error()}, err
	}

	absPath, err := filepath.Abs(opts.ScriptPath)
	if err != nil {
		err = fmt.Errorf("resolve script path: %w", err)
		return RunPlanResult{Error: err.Error()}, err
	}

	opts.PlaywrightVersion = applyResolvedPlaywrightVersion(opts.PlaywrightVersion, opts.Stderr)

	cacheDir, err := VersionedCacheDir(opts.CacheDir, opts.PlaywrightVersion)
	if err != nil {
		return RunPlanResult{Error: err.Error()}, err
	}

	if opts.SkipEnsure {
		if err := os.MkdirAll(cacheDir, 0o755); err != nil {
			err = fmt.Errorf("create cache dir: %w", err)
			return RunPlanResult{Error: err.Error()}, err
		}
	} else {
		cacheDir, err = EnsurePlaywrightVersion(opts.CacheDir, opts.PlaywrightVersion, opts.Stdout, opts.Stderr)
		if err != nil {
			return RunPlanResult{Error: err.Error()}, err
		}
	}

	bootstrapPath := filepath.Join(cacheDir, "bootstrap.cjs")
	if err := os.WriteFile(bootstrapPath, []byte(bootstrapScript), 0o644); err != nil {
		err = fmt.Errorf("write bootstrap.cjs: %w", err)
		return RunPlanResult{Error: err.Error()}, err
	}

	scriptArgs := append([]string(nil), opts.ScriptArgs...)
	nodeArgv := append([]string{"bootstrap.cjs", absPath}, scriptArgs...)

	return RunPlanResult{
		ScriptPath:        absPath,
		ScriptArgs:        scriptArgs,
		CacheDir:          cacheDir,
		PlaywrightVersion: opts.PlaywrightVersion,
		NodeArgv:          nodeArgv,
		BootstrapWritten:  true,
	}, nil
}

// RunFile executes a Playwright .js script via the embedded bootstrap runner.
func RunFile(ctx context.Context, opts RunOptions) error {
	plan, err := BuildRunPlan(opts)
	if err != nil {
		return err
	}

	stdout := opts.Stdout
	if stdout == nil {
		stdout = os.Stdout
	}
	stderr := opts.Stderr
	if stderr == nil {
		stderr = os.Stderr
	}

	bootstrapPath := filepath.Join(plan.CacheDir, "bootstrap.cjs")
	cmd := exec.CommandContext(ctx, "node", append([]string{bootstrapPath, plan.ScriptPath}, plan.ScriptArgs...)...)
	cmd.Dir = plan.CacheDir
	cmd.Env = opts.Launch.ApplyEnv(PlaywrightEnv(plan.CacheDir))
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return &ExitError{Code: exitErr.ExitCode(), Err: fmt.Errorf("playwright script failed: %w", err)}
		}
		return fmt.Errorf("playwright script failed: %w", err)
	}
	return nil
}

// ExitError preserves a subprocess exit code for CLI wrappers.
type ExitError struct {
	Code int
	Err  error
}

func (e *ExitError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return fmt.Sprintf("exit status %d", e.Code)
}

func (e *ExitError) Unwrap() error {
	return e.Err
}

// EnsurePlaywright initializes the npm package cache and installs playwright when needed.
func EnsurePlaywright(cacheDir string, stdout, stderr io.Writer) (string, error) {
	return EnsurePlaywrightVersion(cacheDir, "", stdout, stderr)
}

// EnsurePlaywrightVersion initializes the npm package cache, enforces an
// explicitly requested package version, repairs a missing Chromium executable
// left by an interrupted browser download, and on Linux installs OS shared
// libraries required to launch Chromium (playwright install-deps).
func EnsurePlaywrightVersion(cacheDir, version string, stdout, stderr io.Writer) (string, error) {
	if stdout == nil {
		stdout = os.Stdout
	}
	if stderr == nil {
		stderr = os.Stderr
	}
	version = applyResolvedPlaywrightVersion(version, stderr)

	var err error
	cacheDir, err = VersionedCacheDir(cacheDir, version)
	if err != nil {
		return "", err
	}

	if playwrightReady(cacheDir, version) {
		return cacheDir, nil
	}

	if err := withEnsureLock(cacheDir, func() error {
		return ensurePlaywrightLocked(cacheDir, version, stdout, stderr)
	}); err != nil {
		return "", err
	}

	return cacheDir, nil
}

type playwrightPackageJSON struct {
	Version string `json:"version"`
}

func installedPlaywrightVersion(cacheDir string) (string, error) {
	data, err := os.ReadFile(filepath.Join(cacheDir, "node_modules", "playwright", "package.json"))
	if err != nil {
		return "", err
	}
	var pkg playwrightPackageJSON
	if err := json.Unmarshal(data, &pkg); err != nil {
		return "", err
	}
	if pkg.Version == "" {
		return "", fmt.Errorf("playwright package.json has no version")
	}
	return pkg.Version, nil
}

func playwrightPackageMatches(cacheDir, version string) bool {
	installed, err := installedPlaywrightVersion(cacheDir)
	if err != nil {
		return false
	}
	return version == "" || installed == version
}

func chromiumExecutableExists(cacheDir string) bool {
	const checkScript = `const fs=require("fs");const{chromium}=require("playwright");const p=chromium.executablePath();process.exit(p&&fs.existsSync(p)?0:1)`
	cmd := exec.Command("node", "-e", checkScript)
	cmd.Dir = cacheDir
	cmd.Env = PlaywrightEnv(cacheDir)
	return cmd.Run() == nil
}

// needsChromiumOSDeps is true on Linux, where Chromium needs apt-installed
// shared libraries (libatk, etc.). Overridable in tests.
var needsChromiumOSDeps = func() bool { return runtime.GOOS == "linux" }

// chromiumHostDepsOK reports whether Chromium can load on this host.
// Non-Linux always OK. Overridable in tests.
var chromiumHostDepsOK = defaultChromiumHostDepsOK

func defaultChromiumHostDepsOK(cacheDir string) bool {
	if !needsChromiumOSDeps() {
		return true
	}
	return smokeChromiumLoads(cacheDir)
}

// smokeChromiumLoads resolves the Playwright Chromium binary and probes it.
// Missing shared libraries typically yield exit 127 from the dynamic linker.
func smokeChromiumLoads(cacheDir string) bool {
	const pathScript = `const {chromium}=require("playwright");const p=chromium.executablePath();if(!p)process.exit(1);process.stdout.write(p)`
	cmd := exec.Command("node", "-e", pathScript)
	cmd.Dir = cacheDir
	cmd.Env = PlaywrightEnv(cacheDir)
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	exe := strings.TrimSpace(string(out))
	if exe == "" {
		return false
	}
	if _, err := os.Stat(exe); err != nil {
		return false
	}
	probe := exec.Command(exe, "--version")
	probe.Dir = cacheDir
	return probe.Run() == nil
}

func chromiumReady(cacheDir string) bool {
	return chromiumExecutableExists(cacheDir) && chromiumHostDepsOK(cacheDir)
}

func playwrightReady(cacheDir, version string) bool {
	return playwrightPackageMatches(cacheDir, version) && chromiumReady(cacheDir)
}

func ensurePlaywrightLocked(cacheDir, version string, stdout, stderr io.Writer) error {
	if playwrightReady(cacheDir, version) {
		return nil
	}

	packageJSON := filepath.Join(cacheDir, "package.json")
	if _, err := os.Stat(packageJSON); os.IsNotExist(err) {
		fmt.Fprintln(stdout, "Initializing playwright cache directory...")
		cmd := exec.Command("npm", "init", "-y")
		cmd.Dir = cacheDir
		cmd.Stdout = stdout
		cmd.Stderr = stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("npm init: %w", err)
		}
	}

	if !playwrightPackageMatches(cacheDir, version) {
		target := "playwright"
		if version != "" {
			target += "@" + version
		}
		fmt.Fprintf(stdout, "Installing %s...\n", target)
		cmd := exec.Command("npm", "install", "--save-exact", target)
		cmd.Dir = cacheDir
		cmd.Stdout = stdout
		cmd.Stderr = stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("npm install %s: %w", target, err)
		}
	}

	if !playwrightPackageMatches(cacheDir, version) {
		installed, err := installedPlaywrightVersion(cacheDir)
		if err != nil {
			return fmt.Errorf("verify installed Playwright package: %w", err)
		}
		return fmt.Errorf("installed Playwright version %q does not match requested version %q", installed, version)
	}

	if !chromiumExecutableExists(cacheDir) {
		label := "playwright"
		if version != "" {
			label += "@" + version
		}
		fmt.Fprintf(stdout, "Installing Chromium browser for %s...\n", label)
		cmd := exec.Command("npm", "exec", "--", "playwright", "install", "chromium")
		cmd.Dir = cacheDir
		cmd.Stdout = stdout
		cmd.Stderr = stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("playwright install chromium for %s: %w", label, err)
		}
	}

	if !chromiumExecutableExists(cacheDir) {
		return fmt.Errorf("Playwright Chromium executable is missing after installation")
	}

	if needsChromiumOSDeps() && !chromiumHostDepsOK(cacheDir) {
		fmt.Fprintln(stdout, "Installing Chromium OS dependencies (playwright install-deps)...")
		cmd := exec.Command("npm", "exec", "--", "playwright", "install-deps", "chromium")
		cmd.Dir = cacheDir
		cmd.Stdout = stdout
		cmd.Stderr = stderr
		env, err := installDepsEnv(cacheDir)
		if err != nil {
			return err
		}
		cmd.Env = env
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("install OS dependencies for Chromium failed (apt may need root/sudo; ensure uses Check-Valid-Until=false and Debian main-only sources to avoid flaky security.debian.org 404s). Fix: cd %s && npm exec -- playwright install-deps chromium: %w", cacheDir, err)
		}
		if !chromiumHostDepsOK(cacheDir) {
			return fmt.Errorf("Chromium still cannot load after install-deps (missing shared libraries?). Fix: cd %s && npm exec -- playwright install-deps chromium", cacheDir)
		}
	}

	return nil
}

func aptInstallDepsConfPath(cacheDir string) string {
	return filepath.Join(cacheDir, ".playwright-debug-apt.conf")
}

func aptMainOnlySourcesPath(cacheDir string) string {
	return filepath.Join(cacheDir, ".playwright-debug-apt-sources.list")
}

func aptEmptySourcePartsDir(cacheDir string) string {
	return filepath.Join(cacheDir, ".playwright-debug-apt-sourceparts")
}

// ensureAptInstallDepsConf writes APT_CONFIG (+ optional Debian main-only
// sources) so install-deps can ignore expired Valid-Until and skip flaky
// security.debian.org CDN 404s.
func ensureAptInstallDepsConf(cacheDir string) (string, error) {
	var b strings.Builder
	b.WriteString("Acquire::Check-Valid-Until \"false\";\n")

	if host := lookupHostOS(); strings.EqualFold(strings.TrimSpace(host.ID), "debian") {
		sourcesPath, err := ensureDebianMainOnlySources(cacheDir, host)
		if err != nil {
			return "", err
		}
		partsDir := aptEmptySourcePartsDir(cacheDir)
		if err := os.MkdirAll(partsDir, 0o755); err != nil {
			return "", fmt.Errorf("create empty apt sourceparts dir: %w", err)
		}
		fmt.Fprintf(&b, "Dir::Etc::sourcelist \"%s\";\n", sourcesPath)
		fmt.Fprintf(&b, "Dir::Etc::sourceparts \"%s\";\n", partsDir)
	}

	path := aptInstallDepsConfPath(cacheDir)
	if err := os.WriteFile(path, []byte(b.String()), 0o644); err != nil {
		return "", fmt.Errorf("write apt conf for install-deps: %w", err)
	}
	return path, nil
}

func ensureDebianMainOnlySources(cacheDir string, host HostOS) (string, error) {
	codename := strings.TrimSpace(host.VersionCodename)
	if codename == "" {
		codename = debianCodenameForVersionID(host.VersionID)
	}
	if codename == "" {
		return "", fmt.Errorf("cannot derive Debian codename for VERSION_ID=%q", host.VersionID)
	}
	// Main + updates only — no security.debian.org (Fastly edges often 404
	// bullseye-security pool files while indexes still list them).
	body := fmt.Sprintf(
		"deb http://deb.debian.org/debian %s main\ndeb http://deb.debian.org/debian %s-updates main\n",
		codename, codename,
	)
	path := aptMainOnlySourcesPath(cacheDir)
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		return "", fmt.Errorf("write main-only apt sources: %w", err)
	}
	return path, nil
}

func debianCodenameForVersionID(versionID string) string {
	switch strings.TrimSpace(versionID) {
	case "11":
		return "bullseye"
	case "12":
		return "bookworm"
	case "13":
		return "trixie"
	default:
		return ""
	}
}

// installDepsEnv is the process env for playwright install-deps: PlaywrightEnv
// plus DEBIAN_FRONTEND and APT_CONFIG (Valid-Until bypass; Debian main-only).
func installDepsEnv(cacheDir string) ([]string, error) {
	aptConf, err := ensureAptInstallDepsConf(cacheDir)
	if err != nil {
		return nil, err
	}
	env := PlaywrightEnv(cacheDir)
	env = unsetEnvKey(env, "APT_CONFIG")
	env = unsetEnvKey(env, "DEBIAN_FRONTEND")
	return append(env,
		"APT_CONFIG="+aptConf,
		"DEBIAN_FRONTEND=noninteractive",
	), nil
}

func unsetEnvKey(env []string, key string) []string {
	prefix := key + "="
	out := make([]string, 0, len(env))
	for _, e := range env {
		if strings.HasPrefix(e, prefix) {
			continue
		}
		out = append(out, e)
	}
	return out
}

// PlaywrightEnv returns process environment with NODE_PATH set for the cache directory.
func PlaywrightEnv(cacheDir string) []string {
	nodePath := filepath.Join(cacheDir, "node_modules")
	env := os.Environ()
	for i, e := range env {
		if strings.HasPrefix(e, "NODE_PATH=") {
			existing := strings.TrimPrefix(e, "NODE_PATH=")
			if existing != "" {
				nodePath = nodePath + string(os.PathListSeparator) + existing
			}
			env = append(env[:i], env[i+1:]...)
			break
		}
	}
	return append(env, "NODE_PATH="+nodePath)
}
