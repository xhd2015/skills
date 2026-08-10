# Scenario

**Feature**: bare file alias forwards trailing script arguments

```
user -> playwright-debug CLI (print-argv.js foo bar) -> ["foo","bar"]
```

## Steps

1. Set bare `print-argv.js` path plus `foo` and `bar` on `req.Args` (no `run` prefix).

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{fixturePath(d, "print-argv.js"), "foo", "bar"}
	return nil
}
```