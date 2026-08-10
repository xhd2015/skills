# Scenario

**Feature**: `go-best-practice skill --install` resolves targets and nested extras

```
# dry-run previews install without writing files
user -> go-best-practice skill --install --dry-run -> [dry-run] stdout

# real install writes SKILL.md and nested cli/skill-cli/TOPIC.md
user -> go-best-practice skill --install -> .agents/skills/go-best-practice/
```

## Preconditions

- `req.Binary` and `req.SkillName` are set by the `go-best-practice` grouping setup.

## Steps

1. Each leaf sets `req.Args` for local, global, or real install behavior.

## Context

- Default local target is `.agents/skills/go-best-practice` relative to the work dir.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	if req.Binary == "" {
		t.Fatal("req.Binary must be set by go-best-practice setup")
	}
	if req.SkillName == "" {
		t.Fatal("req.SkillName must be set by go-best-practice setup")
	}
	return nil
}
```
