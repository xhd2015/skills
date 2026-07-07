# Scenario

**Feature**: bare existing `.js` path is file-mode alias

```
user -> playwright-debug CLI (<simple-eval.js>) -> file-mode-ok
```

## Steps

1. Each descendant sets the bare file path on `req.Args`.

```go
func Setup(t *testing.T, req *Request) error {
	if req.Args == nil {
		req.Args = []string{}
	}
	return nil
}
```