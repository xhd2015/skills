## Preconditions
- The test is NOT inside a git repository.
- No origin remote exists to auto-detect.

## Steps
1. Do NOT create a git repo.
2. Run `ci --logs --workflow test`.
3. Check that the command fails with an error about not being in a git repo.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"--logs", "--workflow", "test"}
	req.InGitRepo = false
	return nil
}
```
