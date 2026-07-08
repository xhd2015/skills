# skills

CLI tools for AI-assisted development workflows.

## Tools

### go-best-practice

Index of Go best-practice recipes: project scaffolding, CLI flag parsing, command execution, and more.

```bash
go install github.com/xhd2015/skills/cmd/go-best-practice@latest
go-best-practice                    # list topics
go-best-practice flags-parsing      # read a topic
go-best-practice install --cursor   # install as cursor skill
```

### playwright-debug

Run Playwright browser automation scripts with automatic setup.

```bash
go install github.com/xhd2015/skills/cmd/playwright-debug@latest
playwright-debug 'await page.goto("https://example.com"); console.log(await page.title());'
playwright-debug skill install --cursor   # install as cursor skill
```
