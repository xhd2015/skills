## Preconditions
- Tests in this branch exercise how the repository is resolved from CLI arguments and git context.

## Context
- When `owner/repo` is omitted, the tool auto-detects from `git remote get-url origin`.
- When not in a git repository and no explicit repo is given, the command must fail.

```go
func Setup(t *testing.T, req *Request) error {
	// Repo-resolution branch: tests run without git auto-detect unless a leaf overrides.
	req.InGitRepo = false
	return nil
}
```