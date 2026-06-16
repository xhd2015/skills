## Steps
1. Run `issues otherowner/otherrepo` without git repo.
2. Expect issues for the explicit repo.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"issues", "otherowner/otherrepo"}
	req.InGitRepo = false
	req.MockIssues = []MockIssueListItem{
		{
			Number: 7, Title: "Explicit issue", State: "open", User: "erin",
			HTMLURL: "https://github.com/otherowner/otherrepo/issues/7",
			Labels: []string{"bug"},
		},
	}
	return nil
}
```