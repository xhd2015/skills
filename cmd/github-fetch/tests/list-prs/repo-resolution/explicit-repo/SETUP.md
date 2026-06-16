## Preconditions
- No git repository is required because `owner/repo` is given explicitly.
- The mock API returns open PRs for `otherowner/otherrepo`.

## Steps
1. Run `prs otherowner/otherrepo` outside a git repo.
2. Verify stdout shows the explicit repo and listed PRs.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"prs", "otherowner/otherrepo"}
	req.InGitRepo = false
	req.MockPRs = []MockPRListItem{
		{
			Number: 10, Title: "Explicit repo PR", State: "open", User: "carol",
			HTMLURL: "https://github.com/otherowner/otherrepo/pull/10",
		},
	}
	return nil
}
```