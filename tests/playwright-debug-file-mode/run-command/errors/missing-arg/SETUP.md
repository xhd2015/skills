# Scenario

**Feature**: `run` without a file argument fails

```
user -> playwright-debug CLI (run) -> error: file required
```

## Steps

1. Set `req.Args = []string{"run"}`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"run"}
	return nil
}
```