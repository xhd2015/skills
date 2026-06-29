# Scenario

**Feature**: batch update with zero installed targets reports each skill

```
no pre-install -> HandleUpdateMany(registry) -> skill not installed lines, no SKILL.md created
```

## Preconditions

- No registry skill has `SKILL.md` at default `.agents/skills/<name>` paths.

## Steps

1. Run `HandleUpdateMany` with default target resolution (no `--global`).

```go
func Setup(t *testing.T, req *Request) error {
	req.PreInstalls = nil
	req.UpdateArgs = nil
	return nil
}
```