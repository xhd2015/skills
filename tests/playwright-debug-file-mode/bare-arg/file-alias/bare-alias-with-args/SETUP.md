# Scenario

**Feature**: bare file alias forwards trailing script arguments

```
user -> playwright-debug CLI (print-argv.js foo bar) -> ["foo","bar"]
```

## Steps

1. Set bare `print-argv.js` path plus `foo` and `bar` on `req.Args` (no `run` prefix).

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{fixturePath("print-argv.js"), "foo", "bar"}
	return nil
}
```