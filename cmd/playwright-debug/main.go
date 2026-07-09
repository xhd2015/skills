package main

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/xhd2015/skills/playwrightdebug"
	"github.com/xhd2015/skills/skillcmd"
)

//go:embed SKILL.md
var skillTemplate string

const skillName = "playwright-debug"

const help = `Usage: playwright-debug <command> [ARGS]

Commands:
  run <file.js> [args...]  Run an existing Playwright .js script file
  skill --show [--header]  Show the content of SKILL.md
  skill --install [<dir>]  Install skill SKILL.md to a directory
  skill --list             Print the skill name

Invocation modes:
  playwright-debug <file.js> [args...]         Run script file (file alias)
  playwright-debug run <file.js> [args...]     Explicit file mode (file required)
  playwright-debug -e '<script>' [args...]     Adhoc eval (short flag)
  playwright-debug --eval '<script>' [args...] Adhoc eval (long flag)
  playwright-debug '<script>'                  Eval when arg is not an existing .js file

The run command requires an existing .js script file on disk.

Trailing arguments after a script file or eval snippet are passed through to
the Node subprocess as process.argv (from index 3).

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
	if len(args) == 1 && (args[0] == "-h" || args[0] == "--help") {
		fmt.Print(help)
		return nil
	}
	if len(args) > 1 && (args[0] == "-h" || args[0] == "--help") {
		return handleRunEval(args[1], args[2:])
	}

	if script, rest, ok := extractEvalFlag(args); ok {
		if script == "" {
			return fmt.Errorf("-e/--eval requires a script argument")
		}
		return handleRunEval(script, rest)
	}

	switch args[0] {
	case "run":
		if len(args) < 2 {
			return fmt.Errorf("run requires a .js script file argument")
		}
		if err := playwrightdebug.ValidateScriptFile(args[1]); err != nil {
			return err
		}
		return handleRunFile(args[1], args[2:]...)
	case "skill":
		return singleSkill().Handle(args[1:])
	default:
		if isScriptFile(args[0]) {
			return handleRunFile(args[0], args[1:]...)
		}
		return handleRunEval(strings.Join(args, " "), nil)
	}
}

func singleSkill() *skillcmd.SingleSkill {
	return &skillcmd.SingleSkill{
		Name:        skillName,
		RootContent: skillTemplate,
		Usage:       "skill --install",
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

func handleRunFile(scriptPath string, scriptArgs ...string) error {
	err := playwrightdebug.RunFile(context.Background(), playwrightdebug.RunOptions{
		ScriptPath: scriptPath,
		ScriptArgs: scriptArgs,
		Stdout:     os.Stdout,
		Stderr:     os.Stderr,
	})
	if exitErr, ok := err.(*playwrightdebug.ExitError); ok {
		os.Exit(exitErr.Code)
	}
	return err
}

func handleRunEval(script string, scriptArgs []string) error {
	dir, err := playwrightdebug.EnsurePlaywright("", os.Stdout, os.Stderr)
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

	cmd := exec.Command("node", append([]string{"-e", wrapper}, scriptArgs...)...)
	cmd.Dir = dir
	cmd.Env = playwrightdebug.PlaywrightEnv(dir)
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
