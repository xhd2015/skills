## Steps
1. Run `prs` with no flags or positional arguments.
2. Expect open PRs from auto-detected repo on page 1.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"prs"}
	return nil
}
```