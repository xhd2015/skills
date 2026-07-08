package playwrightdebug

import (
	"context"
	_ "embed"
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
	ScriptPath       string   `json:"script_path"`
	ScriptArgs       []string `json:"script_args"`
	CacheDir         string   `json:"cache_dir"`
	NodeArgv         []string `json:"node_argv"`
	BootstrapWritten bool     `json:"bootstrap_written"`
	Error            string   `json:"error,omitempty"`
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

	cacheDir := opts.CacheDir
	if cacheDir == "" {
		cacheDir = DefaultCacheDir()
	}

	if opts.SkipEnsure {
		if err := os.MkdirAll(cacheDir, 0o755); err != nil {
			err = fmt.Errorf("create cache dir: %w", err)
			return RunPlanResult{Error: err.Error()}, err
		}
	} else {
		cacheDir, err = EnsurePlaywright(cacheDir, opts.Stdout, opts.Stderr)
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
		ScriptPath:       absPath,
		ScriptArgs:       scriptArgs,
		CacheDir:         cacheDir,
		NodeArgv:         nodeArgv,
		BootstrapWritten: true,
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
	cmd.Env = PlaywrightEnv(plan.CacheDir)
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
	if cacheDir == "" {
		cacheDir = DefaultCacheDir()
	}
	if stdout == nil {
		stdout = os.Stdout
	}
	if stderr == nil {
		stderr = os.Stderr
	}

	if playwrightInstalled(cacheDir) {
		return cacheDir, nil
	}

	if err := withEnsureLock(cacheDir, func() error {
		return ensurePlaywrightLocked(cacheDir, stdout, stderr)
	}); err != nil {
		return "", err
	}

	return cacheDir, nil
}

func playwrightInstalled(cacheDir string) bool {
	nodeModules := filepath.Join(cacheDir, "node_modules", "playwright")
	_, err := os.Stat(nodeModules)
	return err == nil
}

func ensurePlaywrightLocked(cacheDir string, stdout, stderr io.Writer) error {
	if playwrightInstalled(cacheDir) {
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

	if playwrightInstalled(cacheDir) {
		return nil
	}

	fmt.Fprintln(stdout, "Installing playwright...")
	cmd := exec.Command("npm", "install", "playwright")
	cmd.Dir = cacheDir
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("npm install playwright: %w", err)
	}

	fmt.Fprintln(stdout, "Installing Chromium browser...")
	npx := "npx"
	if runtime.GOOS == "windows" {
		npx = "npx.cmd"
	}
	cmd = exec.Command(npx, "playwright", "install", "chromium")
	cmd.Dir = cacheDir
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("playwright install chromium: %w", err)
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