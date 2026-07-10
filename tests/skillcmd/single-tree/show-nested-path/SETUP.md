# Scenario

**Feature**: `--show a/b` prints nested TOPIC.md content

```
caller -> SingleSkill.Handle(--show a/b) -> Nested A/B Body from a/b/TOPIC.md
```

## Preconditions

- Tree SingleSkill configured by parent with `a/b/TOPIC.md`.

## Steps

1. Configure Args for nested show.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"--show", "a/b"}
	return nil
}
```
