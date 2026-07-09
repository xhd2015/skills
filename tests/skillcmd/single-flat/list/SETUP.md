# Scenario

**Feature**: `--list` prints the skill Name

```
caller -> SingleSkill.Handle(--list) -> Name
```

## Preconditions

- Flat SingleSkill configured by parent.

## Steps

1. Set Args for this action.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"--list"}
	return nil
}
```
