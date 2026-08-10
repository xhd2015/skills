## Preconditions
- Both `prs` and `pr --list` invoke the same list handler and produce identical output.

```go
var defaultOpenPRs = []MockPRListItem{
	{
		Number: 42, Title: "Fix login redirect", State: "open", User: "alice",
		HTMLURL: "https://github.com/testowner/testrepo/pull/42",
	},
	{
		Number: 41, Title: "Add pagination to API client", State: "open", User: "bob",
		HTMLURL: "https://github.com/testowner/testrepo/pull/41",
	},
}

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.MockPRs = defaultOpenPRs
	return nil
}
```