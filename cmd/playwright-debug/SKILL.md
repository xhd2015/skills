---
name: playwright-debug
description: >-
  Run Playwright scripts for browser automation and frontend debugging.
  Use when the user wants to test a web page, debug UI behavior, monitor
  network requests, or automate browser interactions.
---

# Playwright Debug Skill

A CLI tool for running Playwright browser automation scripts with automatic setup.

## Launch flags (tool-level)

Peeled before the run mode so they work with `run`, file alias, and `-e`:

| Flag | Meaning |
|------|---------|
| `--extension <dir>` | Load unpacked extension (repeatable). Uses `launchPersistentContext` + Chromium `--load-extension`. Defaults to **headed**. |
| `--load-extension <dir>` | Alias of `--extension` |
| `--user-data-dir <dir>` | Chromium profile directory. Works **with or without** `--extension`. Cookies, localStorage, and login state persist across runs and reboots if the path is durable (e.g. under `$HOME`). Never auto-deleted. Without this flag, extension mode uses a temp profile; default mode uses ephemeral `chromium.launch`. |
| `--headed` | Force visible browser |
| `--headless` | Force headless (extensions often fail to load) |

```bash
# Persistent login (survives reboot) — pick a stable path under home/cache
playwright-debug --user-data-dir ~/.cache/pw-debug/my-site -e '
  await page.goto("https://example.com");
  // first run: log in interactively with --headed; later runs reuse cookies
  console.log("cookies", (await context.cookies()).length, "context", !!context);
'

# Same profile + headed for manual login once
playwright-debug --user-data-dir ~/.cache/pw-debug/my-site --headed -e '
  await page.goto("https://example.com/login");
  await page.waitForTimeout(120000); // log in by hand, then Ctrl+C or timeout
'

# Load an unpacked MV3 extension and open a page
playwright-debug --extension /path/to/unpacked-ext -e '
  console.log("context", !!context, "exts", extensionPaths);
  const sw = context.serviceWorkers()[0]
    || await context.waitForEvent("serviceworker", { timeout: 15000 });
  console.log("sw", sw && sw.url());
  await page.goto("http://127.0.0.1:43761/go?session=demo");
  await page.waitForTimeout(2000);
  console.log(await page.title());
'

# Extension + durable profile (login + extension both survive reboot)
playwright-debug --extension ./my-ext --user-data-dir ~/.cache/pw-debug/with-ext run ./check-ext.js
```

Persistent / extension mode injects:

- `context` — Playwright `BrowserContext` (when `--user-data-dir` or `--extension` is set)
- `extensionPaths` — absolute paths of loaded unpacked dirs (or `[]`)
- `browser` — often **`null`** under persistent context; prefer `context`
- `page` — first page in the context

Default (no `--extension`, no `--user-data-dir`) stays headless `chromium.launch`;
`context` is `null` and `extensionPaths` is `[]`.

## Commands

### run — Run a Playwright script file

Runs an existing `.js` script file in Chromium (headless by default). The file is
executed with top-level `await` support and a `require()` resolver relative to
the script directory.

Available variables in file mode:
- `browser` — Chromium browser instance (may be null with `--extension` / `--user-data-dir`)
- `page` — A new browser page
- `chromium` — Playwright's chromium launcher
- `context` — BrowserContext when `--extension` or `--user-data-dir` is set; otherwise null
- `extensionPaths` — string[] of loaded unpacked extension paths
- `require` — Node require scoped to the script file
- `__filename` — Absolute path to the script file
- `__dirname` — Directory containing the script file

```bash
# Run a script file explicitly
playwright-debug run ./my-script.js

# Bare file alias (existing .js path)
playwright-debug ./my-script.js

# Forward trailing args to process.argv (from index 3)
playwright-debug run ./my-script.js -o /tmp/screenshot.png
playwright-debug ./my-script.js --verbose
```

Trailing arguments after a script file or eval snippet are passed through to
the Node subprocess as `process.argv` (from index 3).

`--help` after a script path is forwarded to the script (CLI help is shown only
for bare `-h`/`--help` with no other arguments).

### Eval — Run an adhoc script snippet

Runs a JavaScript snippet in Chromium (headless by default). The script is
automatically wrapped with browser/page setup.

Available variables in eval mode:
- `browser` — Chromium browser instance (may be null with `--extension` / `--user-data-dir`)
- `page` — A new browser page
- `chromium` — Playwright's chromium launcher
- `context` — BrowserContext when `--extension` or `--user-data-dir` is set; otherwise null
- `extensionPaths` — string[] of loaded unpacked extension paths

```bash
# Explicit eval flag
playwright-debug -e 'console.log("eval-ok")'
playwright-debug --eval 'await page.goto("https://example.com"); console.log(await page.title());'

# Forward trailing args to process.argv (from index 3)
playwright-debug -e 'console.log(process.argv.slice(3))' foo bar
playwright-debug --eval 'console.log(process.argv.slice(3))' --output /tmp/out.png

# Bare script string (when the argument is not an existing .js file)
playwright-debug 'await page.goto("https://example.com"); console.log(await page.title());'

# Monitor network requests
playwright-debug '
  page.on("request", req => console.log("[REQ]", req.method(), req.url()));
  page.on("requestfailed", req => console.log("[FAIL]", req.url(), req.failure()?.errorText));
  await page.goto("https://example.com");
  await page.waitForTimeout(5000);
'

# Fill a form and submit
playwright-debug '
  await page.goto("https://example.com/login");
  await page.fill("input[name=email]", "test@example.com");
  await page.fill("input[name=password]", "secret");
  await page.click("button[type=submit]");
  await page.waitForURL("**/dashboard");
  console.log("Logged in:", await page.title());
'
```