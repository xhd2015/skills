# Scenario

**Feature**: `--dir` and positional `<dir>` together error

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"--dir", "vendor/skills", "other"}
	return nil
}
```
