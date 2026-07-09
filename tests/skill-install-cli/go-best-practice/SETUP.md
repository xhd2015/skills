# Scenario

**Feature**: `go-best-practice` skill actions use flag form via skillcmd

```
# go-best-practice routes skill flag actions
user -> go-best-practice skill --show|--install|--list -> skillcmd SingleSkill

# top-level install remains a backward-compatible alias
user -> go-best-practice install -> skill --install equivalent
```

## Preconditions

- `go-best-practice` binary is built from `cmd/go-best-practice`.

## Steps

1. Resolve `req.Binary` via session cache build.
2. Set `req.SkillName = "go-best-practice"`.
3. Leaves configure args and scope for their scenario.

## Context

- Embedded assets include root SKILL.md and nested `path/SKILL.md` extras
  (migrated from `topics/*.md`).

```go
func Setup(t *testing.T, req *Request) error {
	bin, err := buildGoBestPracticeOnce(t)
	if err != nil {
		return err
	}
	req.Binary = bin
	req.SkillName = "go-best-practice"
	return nil
}
```
