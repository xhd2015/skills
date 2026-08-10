# Scenario

**Feature**: nested lib uses explicit page parameter

```
user -> playwright-debug CLI (run explicit-page/main.js) -> explicit-page-ok
```

## Steps

1. Run `testdata/explicit-page/main.js` which passes injected `page` to `./lib/use-page`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"run", fixturePath(d, "explicit-page", "main.js")}
	return nil
}
```