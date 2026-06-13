## Preconditions
- A git repo exists with origin `git@github.com:testowner/testrepo.git`.
- The mock API returns a "test" workflow run with failed status and associated logs.

## Steps
1. Create a git repo with origin `git@github.com:testowner/testrepo.git`.
2. Configure mock API: one run named "test" with status "completed", conclusion "failure", and log content "build step failed".
3. Run `ci --logs --workflow test`.
4. Check that stdout contains a status header showing "test" workflow with "failure", the log content, and the detected GitHub URL.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"--logs", "--workflow", "test"}
	req.OriginURL = "git@github.com:testowner/testrepo.git"
	req.InGitRepo = true
	req.MockRuns = []MockWorkflowRun{
		{
			ID: 12345, Name: "test", Status: "completed", Conclusion: "failure",
			HTMLURL: "https://github.com/testowner/testrepo/actions/runs/12345",
			HeadBranch: "main", Event: "push",
		},
		{
			ID: 12346, Name: "lint", Status: "completed", Conclusion: "success",
			HTMLURL: "https://github.com/testowner/testrepo/actions/runs/12346",
			HeadBranch: "main", Event: "push",
		},
	}
	req.MockJobs = []MockWorkflowJob{
		{ID: 100, Name: "test-job", Conclusion: "failure"},
	}
	req.MockLogs = "build step failed\nerror: test failed\n"
	return nil
}
```
