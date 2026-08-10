# Scenario

**Feature**: bare non-file string is eval shorthand

```
user -> playwright-debug CLI ('console.log("eval-ok")') -> eval-ok
```

## Steps

1. Descendant sets bare eval script on `req.Args`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	if req.Args == nil {
		req.Args = []string{}
	}
	return nil
}
```