## Preconditions
- A git repo exists with origin `git@github.com:testowner/testrepo.git`.
- The mock API returns multiple workflow runs.
- `--logs` is set but no `--workflow` filter.
- The latest run should be shown.

## Steps
1. Create a git repo with origin `git@github.com:testowner/testrepo.git`.
2. Configure mock API with two runs: "test" (first, success) and "lint" (second, failed).
3. Run `ci --logs`.
4. Check that stdout contains a status header for the latest run ("lint", failed) and its logs.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"--logs"}
	req.OriginURL = "git@github.com:testowner/testrepo.git"
	req.InGitRepo = true
	req.MockRuns = []MockWorkflowRun{
		{
			ID: 100, Name: "test", Status: "completed", Conclusion: "success",
			HTMLURL: "https://github.com/testowner/testrepo/actions/runs/100",
			HeadBranch: "main", Event: "push",
		},
		{
			ID: 200, Name: "lint", Status: "completed", Conclusion: "failure",
			HTMLURL: "https://github.com/testowner/testrepo/actions/runs/200",
			HeadBranch: "main", Event: "push",
		},
	}
	req.MockJobs = []MockWorkflowJob{
		{ID: 200, Name: "lint-job", Conclusion: "failure"},
	}
	req.MockLogs = "linter output\nerror: style violation\n"
	return nil
}
```
