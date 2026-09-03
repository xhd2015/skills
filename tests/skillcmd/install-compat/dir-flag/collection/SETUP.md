# Scenario

**Feature**: `--dir` whose basename is `skills` nests under `<dir>/<skill>`

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"--dir", "vendor/skills"}
	return nil
}
```
