# Scenario

**Feature**: `skill --show skill-cli` prints nested skill-cli topic

```
# flag before path
user -> go-best-practice skill --show skill-cli -> skill-cli nested body
```

## Steps

1. Set Args for flag-before-path order.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"skill", "--show", "skill-cli"}
	return nil
}
```
