## Preconditions
- A git repo exists with origin `git@github.com:testowner/testrepo.git`.
- The mock API returns workflow runs, but none match the `--workflow` filter.
- The mock API also returns a list of workflow files for the helpful error message.

## Steps
1. Create a git repo with origin `git@github.com:testowner/testrepo.git`.
2. Configure mock API with runs named "lint" and "build" only.
3. Configure mock workflow files: `lint.yml`, `build.yml`.
4. Run `ci --workflow test`.
5. Check that the command fails with an error listing available workflows.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"--workflow", "test"}
	req.OriginURL = "git@github.com:testowner/testrepo.git"
	req.InGitRepo = true
	req.MockRuns = []MockWorkflowRun{
		{
			ID: 100, Name: "lint", Status: "completed", Conclusion: "success",
			HTMLURL: "https://github.com/testowner/testrepo/actions/runs/100",
			HeadBranch: "main", Event: "push",
		},
		{
			ID: 200, Name: "build", Status: "completed", Conclusion: "failure",
			HTMLURL: "https://github.com/testowner/testrepo/actions/runs/200",
			HeadBranch: "main", Event: "push",
		},
	}
	req.MockWorkflowFiles = []string{"lint.yml", "build.yml"}
	return nil
}
```
