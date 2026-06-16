## Preconditions
- Pagination is tested with multiple PRs and `--page` flag.
- The mock honors `page` and `per_page` query parameters.

```go
func makePagedPRs(n int) []MockPRListItem {
	prs := make([]MockPRListItem, n)
	for i := 0; i < n; i++ {
		num := n - i
		prs[i] = MockPRListItem{
			Number: num, Title: "PR number " + strconv.Itoa(num), State: "open",
			User: "dev", HTMLURL: "https://github.com/testowner/testrepo/pull/" + strconv.Itoa(num),
		}
	}
	return prs
}

func Setup(t *testing.T, req *Request) error {
	// Default pagination fixture: four open PRs when a leaf does not override MockPRs.
	if len(req.MockPRs) == 0 {
		req.MockPRs = makePagedPRs(4)
	}
	return nil
}
```