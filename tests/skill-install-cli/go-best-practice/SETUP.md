# Scenario

**Feature**: `go-best-practice` gains `skill install` parity with other repo CLIs

```
# go-best-practice routes skill sub-commands
user -> go-best-practice -> handleSkill -> skill install | skill show

# top-level install remains a backward-compatible alias
user -> go-best-practice install -> install.HandleInstall
```

## Preconditions

- `go-best-practice` binary is built from `cmd/go-best-practice`.

## Steps

1. Resolve `req.Binary` via session cache build.
2. Set `req.SkillName = "go-best-practice"`.
3. Leaves configure args and scope for their scenario.

## Context

- Embedded assets include SKILL.md and all `topics/**/*.md` files.

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