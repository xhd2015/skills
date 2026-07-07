# Scenario

**Feature**: `-e` eval prints marker

```
user -> playwright-debug CLI (-e 'console.log("eval-ok")') -> eval-ok
```

## Steps

1. Set `req.Args = []string{"-e", `console.log("eval-ok")`}`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"-e", `console.log("eval-ok")`}
	return nil
}
```