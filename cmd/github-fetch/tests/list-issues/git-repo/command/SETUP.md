## Preconditions
- Mock includes one real issue and one PR (with `pull_request` field); only the issue appears in output.

```go
var defaultOpenIssues = []MockIssueListItem{
	{
		Number: 15, Title: "Login page broken", State: "open", User: "alice",
		HTMLURL: "https://github.com/testowner/testrepo/issues/15",
		Labels: []string{"bug"},
	},
	{
		Number: 42, Title: "Fix login redirect", State: "open", User: "alice",
		HTMLURL: "https://github.com/testowner/testrepo/pull/42",
		IsPullRequest: true,
	},
	{
		Number: 14, Title: "Add dark mode", State: "open", User: "bob",
		HTMLURL: "https://github.com/testowner/testrepo/issues/14",
		Labels: []string{"enhancement"},
	},
}

func Setup(t *testing.T, req *Request) error {
	req.MockIssues = defaultOpenIssues
	return nil
}
```