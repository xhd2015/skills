# Scenario

**Feature**: top-level await works in file mode

```
user -> playwright-debug CLI (run toplevel-await.js) -> toplevel-await-ok
```

## Steps

1. Run `testdata/toplevel-await.js` which uses top-level `await page.goto(...)`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"run", fixturePath("toplevel-await.js")}
	return nil
}
```