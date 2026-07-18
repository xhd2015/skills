package main

import (
	"context"
	_ "embed"
	"encoding/json"
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

const help = `Usage: playwright-debug [launch-flags] <command> [ARGS]

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

Launch flags (tool-level; peeled before script args):
  --extension <dir>        Load unpacked extension (repeatable). Uses
                           launchPersistentContext; defaults to headed.
  --load-extension <dir>   Alias of --extension
  --user-data-dir <dir>    Chromium profile dir (persistent cookies/login).
                           Works with or without --extension. Use a durable
                           path (e.g. ~/.cache/pw-debug/my-site) so login
                           survives reboots. Never auto-deleted.
                           Extension mode without this flag uses a temp profile.
  --headed                 Force visible browser
  --headless               Force headless (extensions may not load)

File mode provides: browser, page, chromium, require, __filename, __dirname,
                    context, extensionPaths
Eval mode provides: browser, page, chromium, context, extensionPaths

  - browser is often null with --extension or --user-data-dir (persistent context)
  - context is set when using --extension or --user-data-dir; null otherwise
  - extensionPaths is string[] of absolute unpacked paths

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

	// Peel launch flags first so they work in any position relative to the mode.
	// For -e/--eval, do not peel inside the script string: extract eval after peel
	// only when flags appear outside the script body (flags before -e are peeled).
	launch, rest, err := playwrightdebug.ExtractLaunchFlags(args)
	if err != nil {
		return err
	}
	args = rest

	if len(args) == 0 {
		fmt.Print(help)
		return nil
	}

	// help with leftover? rare
	if len(args) == 1 && (args[0] == "-h" || args[0] == "--help") {
		fmt.Print(help)
		return nil
	}
	if len(args) > 1 && (args[0] == "-h" || args[0] == "--help") {
		return handleRunEval(args[1], args[2:], launch)
	}

	if script, rest, ok := extractEvalFlag(args); ok {
		if script == "" {
			return fmt.Errorf("-e/--eval requires a script argument")
		}
		return handleRunEval(script, rest, launch)
	}

	switch args[0] {
	case "run":
		if len(args) < 2 {
			return fmt.Errorf("run requires a .js script file argument")
		}
		if err := playwrightdebug.ValidateScriptFile(args[1]); err != nil {
			return err
		}
		return handleRunFile(args[1], args[2:], launch)
	case "skill":
		return singleSkill().Handle(args[1:])
	default:
		if isScriptFile(args[0]) {
			return handleRunFile(args[0], args[1:], launch)
		}
		return handleRunEval(strings.Join(args, " "), nil, launch)
	}
}

func singleSkill() *skillcmd.SingleSkill {
	return &skillcmd.SingleSkill{
		Name:        skillName,
		RootContent: skillTemplate,
		Usage:       "playwright-debug skill --install",
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

func handleRunFile(scriptPath string, scriptArgs []string, launch playwrightdebug.LaunchOptions) error {
	err := playwrightdebug.RunFile(context.Background(), playwrightdebug.RunOptions{
		ScriptPath: scriptPath,
		ScriptArgs: scriptArgs,
		Stdout:     os.Stdout,
		Stderr:     os.Stderr,
		Launch:     launch,
	})
	if exitErr, ok := err.(*playwrightdebug.ExitError); ok {
		os.Exit(exitErr.Code)
	}
	return err
}

func handleRunEval(script string, scriptArgs []string, launch playwrightdebug.LaunchOptions) error {
	dir, err := playwrightdebug.EnsurePlaywright("", os.Stdout, os.Stderr)
	if err != nil {
		return err
	}

	// Inject trailing script args into process.argv so user scripts that use
	// process.argv.slice(3) see them even when the -e wrapper is multi-line
	// (some Node builds / arg packing edge cases drop post -e args).
	argsJSON, err := json.Marshal(scriptArgs)
	if err != nil {
		return fmt.Errorf("marshal script args: %w", err)
	}

	// Shared launch helper is embedded in bootstrap; for eval, inline a compact
	// duplicate that reads the same env keys (keep in sync with bootstrap.cjs).
	wrapper := fmt.Sprintf(`
const { chromium } = require('playwright');
const fs = require('fs');
const os = require('os');
const path = require('path');

// Ensure process.argv matches: [node, -e, <wrapper>, ...scriptArgs]
const __pwDebugScriptArgs = %s;
process.argv.length = 3;
for (const a of __pwDebugScriptArgs) process.argv.push(a);

function envFlag(name, fallback) {
  const v = process.env[name];
  if (v == null || v === '') return fallback;
  return v;
}

function resolveUserDataDir(allowTemp) {
  let userDataDir = envFlag('PLAYWRIGHT_DEBUG_USER_DATA_DIR', '');
  let createdTemp = false;
  if (userDataDir) {
    fs.mkdirSync(userDataDir, { recursive: true });
    return { userDataDir, createdTemp: false };
  }
  if (allowTemp) {
    userDataDir = fs.mkdtempSync(path.join(os.tmpdir(), 'pw-debug-profile-'));
    createdTemp = true;
  }
  return { userDataDir, createdTemp };
}

async function launchBrowser() {
  const mode = envFlag('PLAYWRIGHT_DEBUG_LAUNCH_MODE', 'default');
  const headed = envFlag('PLAYWRIGHT_DEBUG_HEADED', mode === 'extension' ? '1' : '0') === '1';
  const raw = envFlag('PLAYWRIGHT_DEBUG_EXTENSION_PATHS', '');
  const extPaths = raw ? raw.split(path.delimiter).map(p => p.trim()).filter(Boolean) : [];
  const wantExtension = mode === 'extension' || extPaths.length > 0;
  const explicitProfile = envFlag('PLAYWRIGHT_DEBUG_USER_DATA_DIR', '') !== '';

  if (wantExtension) {
    if (!extPaths.length) throw new Error('extension mode requires PLAYWRIGHT_DEBUG_EXTENSION_PATHS');
    const { userDataDir, createdTemp } = resolveUserDataDir(true);
    const joined = extPaths.join(',');
    const context = await chromium.launchPersistentContext(userDataDir, {
      headless: !headed,
      args: [
        '--disable-extensions-except=' + joined,
        '--load-extension=' + joined,
        '--no-first-run',
        '--no-default-browser-check',
      ],
    });
    const page = context.pages()[0] || await context.newPage();
    const browser = typeof context.browser === 'function' ? context.browser() : null;
    return {
      browser, context, page, extensionPaths: extPaths,
      async close() {
        await context.close();
        if (createdTemp) { try { fs.rmSync(userDataDir, { recursive: true, force: true }); } catch (_) {} }
      },
    };
  }

  if (explicitProfile) {
    const { userDataDir } = resolveUserDataDir(false);
    const context = await chromium.launchPersistentContext(userDataDir, {
      headless: !headed,
      args: ['--no-first-run', '--no-default-browser-check'],
    });
    const page = context.pages()[0] || await context.newPage();
    const browser = typeof context.browser === 'function' ? context.browser() : null;
    return {
      browser, context, page, extensionPaths: [],
      async close() { await context.close(); },
    };
  }

  const browser = await chromium.launch({ headless: !headed });
  const page = await browser.newPage();
  return {
    browser, context: null, page, extensionPaths: [],
    async close() { await browser.close(); },
  };
}

(async () => {
  const launched = await launchBrowser();
  const browser = launched.browser;
  const page = launched.page;
  const context = launched.context;
  const extensionPaths = launched.extensionPaths;
  try {
    %s
  } finally {
    await launched.close();
  }
})();
`, string(argsJSON), script)

	cmd := exec.Command("node", append([]string{"-e", wrapper}, scriptArgs...)...)
	cmd.Dir = dir
	cmd.Env = launch.ApplyEnv(playwrightdebug.PlaywrightEnv(dir))
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
