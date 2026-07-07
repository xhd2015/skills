# Scenario

**Feature**: `run` rejects more than one file argument

```
user -> playwright-debug CLI (run a.js b.js) -> exactly one file error
```

## Steps

1. Set `req.Args` with two existing fixture paths.

```go
func Setup(t *testing.T, req *Request) error {
	a := fixturePath("simple-eval.js")
	b := fixturePath("toplevel-await.js")
	req.Args = []string{"run", a, b}
	return nil
}
```