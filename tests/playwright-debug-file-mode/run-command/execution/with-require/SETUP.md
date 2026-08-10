# Scenario

**Feature**: relative require() resolves from script directory

```
user -> playwright-debug CLI (run with-require/main.js) -> require-ok
```

## Steps

1. Run `testdata/with-require/main.js` which requires `./lib/helper`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"run", fixturePath(d, "with-require", "main.js")}
	return nil
}
```