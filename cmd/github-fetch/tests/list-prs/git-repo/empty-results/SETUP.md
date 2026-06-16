## Steps
1. Configure mock with no PRs.
2. Run `prs` in git repo.
3. Expect success with an empty-results message.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"prs"}
	req.MockPRs = []MockPRListItem{}
	return nil
}
```