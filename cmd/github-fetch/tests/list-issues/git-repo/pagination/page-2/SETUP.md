## Steps
1. Run `issues --page 2 --per-page 2` with 4 mock issues.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"issues", "--page", "2", "--per-page", "2"}
	req.MockIssues = makePagedIssues(4)
	return nil
}
```