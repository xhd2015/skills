# Scenario

**Feature**: help text documents the status sub-command

```
# help surfaces status in CLI docs
github-fetch -h -> command list includes status
github-fetch status -h -> status-specific usage
```

## Preconditions

- `github-fetch` binary is built from this module.
- Help invocations do not require mock API responses (mock server still started by `Run`).

## Steps

1. Leaves set `req.Args` to `["-h"]` or `["status", "-h"]`.
2. Assert help text content.

## Context

- Help tests exit 0 and print usage to stdout.

```go
func Setup(t *testing.T, req *Request) error {
	// Clear root default `status` argv; help leaves set explicit `-h` invocations.
	if len(req.Args) == 1 && req.Args[0] == "status" {
		req.Args = nil
	}
	req.MockAPIFail = false
	return nil
}
```