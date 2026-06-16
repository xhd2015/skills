## Steps
1. Run `issues` with empty mock data.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"issues"}
	req.MockIssues = []MockIssueListItem{}
	return nil
}
```