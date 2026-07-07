# Scenario

**Feature**: simple file script prints marker

```
user -> playwright-debug CLI (run simple-eval.js) -> file-mode-ok
```

## Steps

1. Run `testdata/simple-eval.js` via explicit `run`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"run", fixturePath("simple-eval.js")}
	return nil
}
```