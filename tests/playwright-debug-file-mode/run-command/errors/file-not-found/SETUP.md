# Scenario

**Feature**: `run` with missing file path fails

```
user -> playwright-debug CLI (run missing.js) -> file not found error
```

## Steps

1. Set `req.Args = []string{"run", "missing.js"}`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"run", "missing.js"}
	return nil
}
```