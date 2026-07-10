# Scenario

**Feature**: `--list` derives topics from nested TOPIC.md paths only

```
TreeFS: a/b/TOPIC.md, skill-cli/TOPIC.md
caller -> SingleSkill.Handle(--list)
  -> demo-skill
  -> a/b
  -> skill-cli
```

## Preconditions

- Parent configures TreeFiles with nested TOPIC.md entries.
- No nested SKILL.md is required for topic discovery.

## Steps

1. Set Args to `--list`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"--list"}
	return nil
}
```
