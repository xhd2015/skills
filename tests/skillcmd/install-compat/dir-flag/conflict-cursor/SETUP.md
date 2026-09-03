# Scenario

**Feature**: `--dir` combined with `--cursor` errors

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"--dir", "vendor/skills", "--cursor"}
	return nil
}
```
