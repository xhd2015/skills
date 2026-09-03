# Scenario

**Feature**: `<topic> --install --dir <collection>` peels topic and nests under collection

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"skill-cli", "--install", "--dir", "vendor/skills", "--dry-run"}
	return nil
}
```
