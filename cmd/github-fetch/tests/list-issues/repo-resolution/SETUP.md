## Preconditions
- Tests in this branch exercise repository resolution for `issues` listing.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.InGitRepo = false
	return nil
}
```