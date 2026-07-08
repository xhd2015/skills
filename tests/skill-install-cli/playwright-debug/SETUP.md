# Scenario

**Feature**: `playwright-debug skill install` global dry-run (regression / parity)

```
# existing skill install support must keep working
user -> playwright-debug skill install --global --dry-run -> HOME/.agents/skills/playwright-debug
```

## Preconditions

- `playwright-debug` binary is built from `cmd/playwright-debug`.

## Steps

1. Resolve `req.Binary` via session cache build.
2. Set `req.SkillName = "playwright-debug"`.
3. Leaf sets global dry-run args.

## Context

- Parity check that global dry-run path resolution matches other repo CLIs.

```go
func Setup(t *testing.T, req *Request) error {
	bin, err := buildPlaywrightDebugOnce(t)
	if err != nil {
		return err
	}
	req.Binary = bin
	req.SkillName = "playwright-debug"
	return nil
}
```