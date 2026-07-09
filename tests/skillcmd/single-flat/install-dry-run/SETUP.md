# Scenario

**Feature**: `--install --dry-run` previews default agents target

```
caller -> SingleSkill.Handle(--install --dry-run) -> [dry-run] .agents/skills/demo-skill
```

## Preconditions

- Flat SingleSkill configured by parent.

## Steps

1. Set Args for this action.

```go
func Setup(t *testing.T, req *Request) error {
	req.UseWorkDir = true
	req.Args = []string{"--install", "--dry-run"}
	return nil
}
```
