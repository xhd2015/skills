## Steps
1. Run `issues --state closed` with one closed and one open issue.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"issues", "--state", "closed"}
	req.MockIssues = []MockIssueListItem{
		{
			Number: 88, Title: "Resolved bug", State: "closed", User: "frank",
			HTMLURL: "https://github.com/testowner/testrepo/issues/88",
		},
		{
			Number: 15, Title: "Still open", State: "open", User: "alice",
			HTMLURL: "https://github.com/testowner/testrepo/issues/15",
		},
	}
	return nil
}
```