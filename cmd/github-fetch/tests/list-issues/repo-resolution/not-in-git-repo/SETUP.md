## Steps
1. Run `issues` outside a git repo with no explicit `owner/repo`.
2. Expect an error about missing git context.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"issues"}
	req.InGitRepo = false
	return nil
}
```