# Scenario

**Feature**: `run` forwards `--help` to the script (not CLI help)

```
user -> playwright-debug CLI (run print-help.js --help) -> SCRIPT_HELP_OK
```

## Steps

1. Run `print-help.js` with `--help` after the file path.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"run", fixturePath("print-help.js"), "--help"}
	return nil
}
```