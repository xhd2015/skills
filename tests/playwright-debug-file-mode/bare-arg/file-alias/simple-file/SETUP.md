# Scenario

**Feature**: file alias runs simple fixture

```
user -> playwright-debug CLI (simple-eval.js) -> file-mode-ok
```

## Steps

1. Set `req.Args = []string{<absolute simple-eval.js>}` without `run` prefix.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{fixturePath("simple-eval.js")}
	return nil
}
```