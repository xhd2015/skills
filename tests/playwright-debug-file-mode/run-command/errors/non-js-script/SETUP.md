# Scenario

**Feature**: `run` with inline script does not eval-fallback

```
user -> playwright-debug CLI (run 'await page.goto("x")') -> requires existing .js file error
```

## Steps

1. Set `req.Args = []string{"run", `await page.goto("x")`}`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"run", `await page.goto("x")`}
	return nil
}
```