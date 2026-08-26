# skills

CLI tools for AI-assisted development workflows.

## Tools

### go-best-practice

**Migrated** to https://github.com/xhd2015/go-best-practice  
(`skills/cmd/go-best-practice` is only a redirect stub for old installs.)

```bash
go install github.com/xhd2015/go-best-practice/cmd/go-best-practice@latest
go-best-practice skill --list
go-best-practice skill --show cli/config
go-best-practice skill --install --global
```

### playwright-debug

Run Playwright browser automation scripts with automatic setup.

```bash
go install github.com/xhd2015/skills/cmd/playwright-debug@latest
playwright-debug 'await page.goto("https://example.com"); console.log(await page.title());'
playwright-debug skill install --cursor   # install as cursor skill
```
