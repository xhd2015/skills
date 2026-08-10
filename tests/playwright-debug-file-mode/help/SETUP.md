# Scenario

**Feature**: help is shown for no args or `-h` / `--help`

```
# empty invocation prints usage
user -> playwright-debug CLI (no args) -> usage on stdout

# help flag anywhere in known positions
user -> playwright-debug CLI (-h) -> usage on stdout
```

## Preconditions

- Root setup built the `playwright-debug` binary.

## Steps

1. Each leaf sets `req.Args` for its help invocation.
2. Assert exit 0 and usage text documents file mode, eval flags, and file alias.

## Context

- Help text must end with a trailing newline after the last content line.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	if req.Args == nil {
		req.Args = []string{}
	}
	return nil
}
```