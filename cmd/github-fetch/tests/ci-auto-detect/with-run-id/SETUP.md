## Preconditions
- A git repo exists with origin `git@github.com:testowner/testrepo.git`.
- The mock API returns multiple workflow runs.
- `--run-id 200` explicitly targets a specific run.

## Steps
1. Create a git repo with origin `git@github.com:testowner/testrepo.git`.
2. Configure mock API with two runs: ID 100 ("test", success) and ID 200 ("lint", failure).
3. Configure mock jobs and logs for the targeted run.
4. Run `ci --run-id 200 --logs`.
5. Check that stdout shows logs for the specified run #200, not run #100.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"--run-id", "200", "--logs"}
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
	req.MockLogs = "linter error: unused variable\n"
	return nil
}
```
