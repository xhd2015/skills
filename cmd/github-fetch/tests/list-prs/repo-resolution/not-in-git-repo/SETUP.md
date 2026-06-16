## Preconditions
- The test is NOT inside a git repository.
- No `owner/repo` positional argument is provided.

## Steps
1. Do NOT create a git repo.
2. Run `prs` with no arguments.
3. Expect an error about missing git repository for auto-detection.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"prs"}
	req.InGitRepo = false
	return nil
}
```