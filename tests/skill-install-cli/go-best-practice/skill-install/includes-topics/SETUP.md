# Scenario

**Feature**: real install copies nested skill-cli/TOPIC.md alongside root SKILL.md

```
# non-dry-run install writes skill tree (Shape 3 nested paths)
user -> go-best-practice skill --install -> .agents/skills/go-best-practice/SKILL.md + skill-cli/TOPIC.md
```

## Preconditions

- Command runs in an isolated work directory with no pre-existing install.

## Steps

1. Set `req.UseWorkDir = true`.
2. Set `req.Args = ["skill", "--install"]` (no `--dry-run`).

```go
func Setup(t *testing.T, req *Request) error {
	req.UseWorkDir = true
	req.Args = []string{"skill", "--install"}
	return nil
}
```
