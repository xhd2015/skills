# Scenario

**Feature**: file alias runs simple fixture

```
user -> playwright-debug CLI (simple-eval.js) -> file-mode-ok
```

## Steps

1. Set `req.Args = []string{<absolute simple-eval.js>}` without `run` prefix.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{fixturePath(d, "simple-eval.js")}
	return nil
}
```