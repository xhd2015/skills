# Scenario

**Feature**: `--show a/b` prints nested SKILL.md content

```
caller -> SingleSkill.Handle(--show a/b) -> Nested A/B Body
```

## Preconditions

- Tree SingleSkill configured by parent.

## Steps

1. Configure Args (and ExtraFiles when installing).

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"--show", "a/b"}
	return nil
}
```
