# Scenario

**Feature**: `run` with missing file path fails

```
user -> playwright-debug CLI (run missing.js) -> file not found error
```

## Steps

1. Set `req.Args = []string{"run", "missing.js"}`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"run", "missing.js"}
	return nil
}
```