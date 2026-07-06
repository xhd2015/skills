# Scenario

**Feature**: `github-fetch status -h` prints status usage

```
# status subcommand help
github-fetch status -h -> status usage text on stdout
```

## Preconditions

- `status` sub-command is registered in the CLI.

## Steps

1. Set `req.Args = []string{"status", "-h"}`.
2. Run and assert status-specific help text.

## Context

- Should mention `status` and read as usage/help output.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"status", "-h"}
	return nil
}
```