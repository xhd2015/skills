## Preconditions
- Tests in this branch exercise repository resolution for `issues` listing.

```go
func Setup(t *testing.T, req *Request) error {
	req.InGitRepo = false
	return nil
}
```