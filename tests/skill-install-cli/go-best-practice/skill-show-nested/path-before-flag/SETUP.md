# Scenario

**Feature**: `skill cli/skill-cli --show` prints the same nested topic

```
# path before flag
user -> go-best-practice skill cli/skill-cli --show -> skill-cli nested body
```

## Steps

1. Set Args for path-before-flag order.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"skill", "cli/skill-cli", "--show"}
	return nil
}
```
