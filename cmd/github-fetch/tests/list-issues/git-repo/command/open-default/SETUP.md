## Steps
1. Run `issues` with no arguments.
2. Expect open issues only (PR #42 excluded).

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"issues"}
	return nil
}
```