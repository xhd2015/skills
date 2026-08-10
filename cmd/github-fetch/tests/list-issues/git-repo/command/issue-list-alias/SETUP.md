## Steps
1. Run `issue --list` alias in git repo.
2. Expect same output as `open-default`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"issue", "--list"}
	return nil
}
```