## Steps
1. Configure mock with more PRs than one page (default per-page 30).
2. Set `HasNextPage` so the mock returns a `Link: rel="next"` header.
3. Run `prs` on page 1.
4. Expect footer hint to use `--page 2`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"prs"}
	req.MockPRs = makePagedPRs(35)
	req.HasNextPage = true
	return nil
}
```