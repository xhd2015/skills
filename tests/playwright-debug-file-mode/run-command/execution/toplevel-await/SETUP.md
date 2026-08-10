# Scenario

**Feature**: top-level await works in file mode

```
user -> playwright-debug CLI (run toplevel-await.js) -> toplevel-await-ok
```

## Steps

1. Run `testdata/toplevel-await.js` which uses top-level `await page.goto(...)`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"run", fixturePath(d, "toplevel-await.js")}
	return nil
}
```