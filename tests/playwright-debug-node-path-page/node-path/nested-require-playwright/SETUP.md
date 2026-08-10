# Scenario

**Feature**: nested require('playwright') resolves via NODE_PATH

```
user -> playwright-debug CLI (run require-playwright-nested/main.js) -> playwright-ok
```

## Steps

1. Run `testdata/require-playwright-nested/main.js` which requires `./lib/check-playwright`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"run", fixturePath(d, "require-playwright-nested", "main.js")}
	return nil
}
```