## Steps
1. Run `pr --list` (alias for `prs`) in the git repo.
2. Expect the same output as `open-default`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"pr", "--list"}
	return nil
}
```