# Scenario

**Feature**: `<topic> --install <collection>` peels topic; positional collection nests

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"skill-cli", "--install", "--dry-run", "vendor/skills"}
	return nil
}
```
