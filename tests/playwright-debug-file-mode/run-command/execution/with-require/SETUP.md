# Scenario

**Feature**: relative require() resolves from script directory

```
user -> playwright-debug CLI (run with-require/main.js) -> require-ok
```

## Steps

1. Run `testdata/with-require/main.js` which requires `./lib/helper`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"run", fixturePath("with-require", "main.js")}
	return nil
}
```