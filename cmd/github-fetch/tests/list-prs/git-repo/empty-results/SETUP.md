## Steps
1. Configure mock with no PRs.
2. Run `prs` in git repo.
3. Expect success with an empty-results message.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"prs"}
	req.MockPRs = []MockPRListItem{}
	return nil
}
```