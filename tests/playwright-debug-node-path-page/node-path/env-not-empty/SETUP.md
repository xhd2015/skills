# Scenario

**Feature**: NODE_PATH env is set when node subprocess runs

```
user -> playwright-debug CLI (run env-check/main.js) -> node-path-set
```

## Steps

1. Run `testdata/env-check/main.js` which prints whether NODE_PATH contains `playwright-debug`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"run", fixturePath("env-check", "main.js")}
	return nil
}
```