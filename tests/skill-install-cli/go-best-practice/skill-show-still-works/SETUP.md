# Scenario

**Feature**: `skill show` continues to work after adding `skill install`

```
# regression: skill show prints embedded SKILL.md
user -> go-best-practice skill show -> stdout contains go-best-practice
```

## Preconditions

- `go-best-practice` binary is built and on `req.Binary`.

## Steps

1. Set `req.Args = ["skill", "show"]`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"skill", "show"}
	return nil
}
```