# Scenario

**Feature**: `run` forwards a single trailing script argument (formerly rejected)

```
user -> playwright-debug CLI (run print-argv.js extra) -> ["extra"]
```

## Steps

1. Run `print-argv.js` with one trailing arg after the file path.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"run", fixturePath(d, "print-argv.js"), "extra"}
	return nil
}
```