# Scenario

**Feature**: bare file alias forwards `--help` to the script (not CLI help)

```
user -> playwright-debug CLI (print-help.js --help) -> SCRIPT_HELP_OK
```

## Steps

1. Set bare `print-help.js` path plus `--help` on `req.Args` (no `run` prefix).

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{fixturePath(d, "print-help.js"), "--help"}
	return nil
}
```