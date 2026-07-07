package main

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/xhd2015/skills/install"
	"github.com/xhd2015/skills/skill_file"
)

//go:embed SKILL.md
var skillTemplate string

//go:embed bootstrap.cjs
var bootstrapScript string

const help = `Usage: playwright-debug <command> [ARGS]

Commands:
  run <file.js>         Run an existing Playwright .js script file
  skill show            Show the content of SKILL.md
  skill install [<dir>] Install skill SKILL.md to a directory

Invocation modes:
  playwright-debug <file.js>              Run script file (file alias)
  playwright-debug run <file.js>          Explicit file mode (file required)
  playwright-debug -e '<script>'          Adhoc eval (short flag)
  playwright-debug --eval '<script>'      Adhoc eval (long flag)
  playwright-debug '<script>'             Eval when arg is not an existing .js file

The run command requires an existing .js script file on disk.

File mode provides: browser, page, chromium, require, __filename, __dirname
Eval mode provides: browser, page, chromium

Options:
  -h, --help    Show this help message
`

func main() {
	if err := handle(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func handle(args []string) error {
	if len(args) == 0 {
		fmt.Print(help)
		return nil
	}

	for _, a := range args {
		if a == "-h" || a == "--help" {
			fmt.Print(help)
			return nil
		}
	}

	if script, rest, ok := extractEvalFlag(args); ok {
		if script == "" {
			return fmt.Errorf("-e/--eval requires a script argument")
		}
		if len(rest) > 0 {
			return fmt.Errorf("unexpected arguments after --eval: %v", rest)
		}
		return handleRunEval(script)
	}

	switch args[0] {
	case "run":
		if len(args) < 2 {
			return fmt.Errorf("run requires a .js script file argument")
		}
		if len(args) > 2 {
			return fmt.Errorf("run accepts exactly one .js file argument")
		}
		if err := validateRunFileArg(args[1]); err != nil {
			return err
		}
		return handleRunFile(args[1])
	case "skill":
		return handleSkill(args[1:])
	default:
		if len(args) == 1 && isScriptFile(args[0]) {
			return handleRunFile(args[0])
		}
		return handleRunEval(strings.Join(args, " "))
	}
}

func extractEvalFlag(args []string) (script string, rest []string, ok bool) {
	for i, a := range args {
		if a == "-e" || a == "--eval" {
			if i+1 >= len(args) {
				return "", nil, true
			}
			return args[i+1], args[i+2:], true
		}
	}
	return "", nil, false
}

func hasJSSuffix(path string) bool {
	return strings.HasSuffix(strings.ToLower(path), ".js")
}

func isScriptFile(path string) bool {
	if !hasJSSuffix(path) {
		return false
	}
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func validateRunFileArg(path string) error {
	if !hasJSSuffix(path) {
		return fmt.Errorf("run requires an existing .js file, not an inline script")
	}
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("script file not found: %s", path)
		}
		return fmt.Errorf("cannot access script file %s: %w", path, err)
	}
	if info.IsDir() {
		return fmt.Errorf("run requires an existing .js file: %s is a directory", path)
	}
	return nil
}

func handleSkill(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("unknown skill sub-command: expected `skill show` or `skill install`")
	}
	switch args[0] {
	case "show":
		rest := args[1:]
		headerOnly := false
		if len(rest) > 0 && rest[0] == "--header" {
			headerOnly = true
			rest = rest[1:]
		}
		if len(rest) > 0 {
			return fmt.Errorf("unknown skill show option: %s", rest[0])
		}
		if headerOnly {
			out, err := skill_file.FormatHeaderWithDelimiters(skillTemplate)
			if err != nil {
				return err
			}
			fmt.Print(out)
			return nil
		}
		fmt.Print(skillTemplate)
		return nil
	case "install":
		return install.HandleInstall(install.InstallOptions{
			SkillDirName: "playwright-debug",
			SkillContent: skillTemplate,
			Usage:        "skill install",
		}, args[1:])
	default:
		return fmt.Errorf("unknown skill sub-command: %s (expected `skill show` or `skill install`)", args[0])
	}
}

func cacheDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".playwright-debug", "node_package")
}

func playwrightEnv(cacheDir string) []string {
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

func ensurePlaywright() (string, error) {
	dir := cacheDir()

	packageJSON := filepath.Join(dir, "package.json")
	if _, err := os.Stat(packageJSON); os.IsNotExist(err) {
		fmt.Println("Initializing playwright cache directory...")
		if err := os.MkdirAll(dir, 0755); err != nil {
			return "", fmt.Errorf("create cache dir: %w", err)
		}
		cmd := exec.Command("npm", "init", "-y")
		cmd.Dir = dir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("npm init: %w", err)
		}
	}

	nodeModules := filepath.Join(dir, "node_modules", "playwright")
	if _, err := os.Stat(nodeModules); os.IsNotExist(err) {
		fmt.Println("Installing playwright...")
		cmd := exec.Command("npm", "install", "playwright")
		cmd.Dir = dir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("npm install playwright: %w", err)
		}

		fmt.Println("Installing Chromium browser...")
		npx := "npx"
		if runtime.GOOS == "windows" {
			npx = "npx.cmd"
		}
		cmd = exec.Command(npx, "playwright", "install", "chromium")
		cmd.Dir = dir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("playwright install chromium: %w", err)
		}
	}

	return dir, nil
}

func handleRunFile(scriptPath string) error {
	if err := validateRunFileArg(scriptPath); err != nil {
		return err
	}

	absPath, err := filepath.Abs(scriptPath)
	if err != nil {
		return fmt.Errorf("resolve script path: %w", err)
	}

	dir, err := ensurePlaywright()
	if err != nil {
		return err
	}

	bootstrapPath := filepath.Join(dir, "bootstrap.cjs")
	if err := os.WriteFile(bootstrapPath, []byte(bootstrapScript), 0644); err != nil {
		return fmt.Errorf("write bootstrap.cjs: %w", err)
	}

	cmd := exec.Command("node", bootstrapPath, absPath)
	cmd.Dir = dir
	cmd.Env = playwrightEnv(dir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		return fmt.Errorf("playwright script failed: %w", err)
	}
	return nil
}

func handleRunEval(script string) error {
	dir, err := ensurePlaywright()
	if err != nil {
		return err
	}

	wrapper := fmt.Sprintf(`
const { chromium } = require('playwright');

(async () => {
  const browser = await chromium.launch({ headless: true });
  const page = await browser.newPage();
  try {
    %s
  } finally {
    await browser.close();
  }
})();
`, script)

	cmd := exec.Command("node", "-e", wrapper)
	cmd.Dir = dir
	cmd.Env = playwrightEnv(dir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		return fmt.Errorf("playwright script failed: %w", err)
	}
	return nil
}