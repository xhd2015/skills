# Scenario

**Feature**: `run` subcommand rejects invalid file arguments before execution

```
user -> playwright-debug CLI (run, bad args) -> routing error on stderr
```

## Preconditions

- No playwright install should be required for these cases.

## Steps

1. Each leaf sets a distinct invalid `run` invocation on `req.Args`.

## Context

- Combined stdout+stderr carries user-facing error messages.

```go
func Setup(t *testing.T, req *Request) error {
	if req.Args == nil {
		req.Args = []string{}
	}
	return nil
}
```