# Scenario

**Feature**: `run` without a file argument fails

```
user -> playwright-debug CLI (run) -> error: file required
```

## Steps

1. Set `req.Args = []string{"run"}`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"run"}
	return nil
}
```