## Preconditions
- Git repo with origin `git@github.com:testowner/testrepo.git` for auto-detection.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.OriginURL = "git@github.com:testowner/testrepo.git"
	req.InGitRepo = true
	return nil
}
```