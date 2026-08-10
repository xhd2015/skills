## Preconditions
- A git repository exists with origin `git@github.com:testowner/testrepo.git`.
- When `owner/repo` is omitted, the tool auto-detects `testowner/testrepo`.

## Steps
1. Create a git repo with the test origin remote.
2. Run list commands without an explicit repo argument.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.OriginURL = "git@github.com:testowner/testrepo.git"
	req.InGitRepo = true
	return nil
}
```