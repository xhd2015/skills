## Preconditions
- A git repo exists WITHOUT an origin remote.
- The auto-detection cannot determine the GitHub repository.

## Steps
1. Create a git repo with initial commit but NO origin remote.
2. Run `ci --logs --workflow test`.
3. Check that the command fails with an error about no origin remote.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"--logs", "--workflow", "test"}
	req.OriginURL = "" // explicitly empty: no origin
	req.InGitRepo = true
	return nil
}
```
