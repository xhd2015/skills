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
// explicitly requested package version, and repairs a missing Chromium
// executable left by an interrupted browser download.
func EnsurePlaywrightVersion(cacheDir, version string, stdout, stderr io.Writer) (string, error) {
	var err error
	cacheDir, err = VersionedCacheDir(cacheDir, version)
	if err != nil {
		return "", err
	}
	if stdout == nil {
		stdout = os.Stdout
	}
	if stderr == nil {
		stderr = os.Stderr
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

func playwrightReady(cacheDir, version string) bool {
	return playwrightPackageMatches(cacheDir, version) && chromiumExecutableExists(cacheDir)
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

	return nil
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
