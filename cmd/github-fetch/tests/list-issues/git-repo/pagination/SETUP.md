## Preconditions
- Pagination tested with multiple issues.

```go
func makePagedIssues(n int) []MockIssueListItem {
	issues := make([]MockIssueListItem, n)
	for i := 0; i < n; i++ {
		num := n - i
		issues[i] = MockIssueListItem{
			Number: num, Title: "Issue number " + strconv.Itoa(num), State: "open",
			User: "dev", HTMLURL: "https://github.com/testowner/testrepo/issues/" + strconv.Itoa(num),
		}
	}
	return issues
}

func Setup(t *testing.T, req *Request) error {
	if len(req.MockIssues) == 0 {
		req.MockIssues = makePagedIssues(4)
	}
	return nil
}
```