## Steps
1. Run `issues` with empty mock data.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"issues"}
	req.MockIssues = []MockIssueListItem{}
	return nil
}
```