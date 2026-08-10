# Scenario

**Feature**: bare eval shorthand prints marker

```
user -> playwright-debug CLI ('console.log("eval-ok")') -> eval-ok
```

## Steps

1. Set `req.Args = []string{`console.log("eval-ok")`}`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{`console.log("eval-ok")`}
	return nil
}
```