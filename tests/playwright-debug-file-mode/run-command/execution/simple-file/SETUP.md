# Scenario

**Feature**: simple file script prints marker

```
user -> playwright-debug CLI (run simple-eval.js) -> file-mode-ok
```

## Steps

1. Run `testdata/simple-eval.js` via explicit `run`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"run", fixturePath(d, "simple-eval.js")}
	return nil
}
```