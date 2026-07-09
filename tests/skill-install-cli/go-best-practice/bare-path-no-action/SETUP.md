# Scenario

**Feature**: bare topic path under skill without action is rejected

```
# missing action flag
user -> go-best-practice skill skill-cli -> error (no bare path action)
```

## Steps

1. Set Args `skill skill-cli` without --show/--install.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"skill", "skill-cli"}
	return nil
}
```
