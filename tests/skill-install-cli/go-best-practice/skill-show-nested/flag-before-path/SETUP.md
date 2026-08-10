# Scenario

**Feature**: `skill --show cli/skill-cli` prints nested skill-cli topic

```
# flag before path
user -> go-best-practice skill --show cli/skill-cli -> skill-cli nested body
```

## Steps

1. Set Args for flag-before-path order.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"skill", "--show", "cli/skill-cli"}
	return nil
}
```
