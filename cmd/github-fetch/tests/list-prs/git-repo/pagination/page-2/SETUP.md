## Steps
1. Configure mock with 4 open PRs.
2. Run `prs --page 2 --per-page 2`.
3. Expect PRs from the second page (numbers 2 and 1).

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"prs", "--page", "2", "--per-page", "2"}
	req.MockPRs = makePagedPRs(4)
	return nil
}
```