## Preconditions
- Git repo with origin `git@github.com:testowner/testrepo.git` for auto-detection.

```go
func Setup(t *testing.T, req *Request) error {
	req.OriginURL = "git@github.com:testowner/testrepo.git"
	req.InGitRepo = true
	return nil
}
```