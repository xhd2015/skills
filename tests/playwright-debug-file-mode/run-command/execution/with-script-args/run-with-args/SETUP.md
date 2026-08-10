# Scenario

**Feature**: `run` forwards multiple trailing script arguments

```
user -> playwright-debug CLI (run print-argv.js -o /tmp/out.png) -> ["-o","/tmp/out.png"]
```

## Steps

1. Run `print-argv.js` with `-o` and `/tmp/out.png` after the file path.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"run", fixturePath(d, "print-argv.js"), "-o", "/tmp/out.png"}
	return nil
}
```