---
name: playwright-debug
description: >-
  Run Playwright scripts for browser automation and frontend debugging.
  Use when the user wants to test a web page, debug UI behavior, monitor
  network requests, or automate browser interactions.
---

# Playwright Debug Skill

A CLI tool for running Playwright browser automation scripts with automatic setup.

## Commands

### run — Run a Playwright script file

Runs an existing `.js` script file in a headless Chromium browser. The file is
executed with top-level `await` support and a `require()` resolver relative to
the script directory.

Available variables in file mode:
- `browser` — Chromium browser instance
- `page` — A new browser page
- `chromium` — Playwright's chromium launcher
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

Runs a JavaScript snippet in a headless Chromium browser. The script is
automatically wrapped with browser/page setup.

Available variables in eval mode:
- `browser` — Chromium browser instance
- `page` — A new browser page
- `chromium` — Playwright's chromium launcher

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