# Scenario

**Feature**: `--dir` whose basename matches the skill name installs directly there

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"--dir", "out/demo-skill"}
	return nil
}
```
