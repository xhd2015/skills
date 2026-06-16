## Steps
1. Configure mock with one open and one closed PR.
2. Run `prs --state closed`.
3. Expect only the closed PR in output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"prs", "--state", "closed"}
	req.MockPRs = []MockPRListItem{
		{
			Number: 99, Title: "Merged feature", State: "closed", User: "dana",
			HTMLURL: "https://github.com/testowner/testrepo/pull/99",
		},
		{
			Number: 42, Title: "Still open", State: "open", User: "alice",
			HTMLURL: "https://github.com/testowner/testrepo/pull/42",
		},
	}
	return nil
}
```