## Preconditions
- A git repo exists with origin `git@github.com:testowner/testrepo.git`.
- The mock API returns a list of recent workflow runs.
- No flags are specified.

## Steps
1. Create a git repo with origin `git@github.com:testowner/testrepo.git`.
2. Configure mock API with two workflow runs: "test" (failed) and "lint" (success).
3. Run `ci` with no arguments.
4. Check that stdout lists both runs with their statuses and the detected GitHub URL.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{}
	req.OriginURL = "git@github.com:testowner/testrepo.git"
	req.InGitRepo = true
	req.MockRuns = []MockWorkflowRun{
		{
			ID: 100, Name: "test", Status: "completed", Conclusion: "failure",
			HTMLURL: "https://github.com/testowner/testrepo/actions/runs/100",
			HeadBranch: "main", Event: "push",
		},
		{
			ID: 200, Name: "lint", Status: "completed", Conclusion: "success",
			HTMLURL: "https://github.com/testowner/testrepo/actions/runs/200",
			HeadBranch: "main", Event: "push",
		},
	}
	req.MockDefaultBranch = "main"
	return nil
}
```
