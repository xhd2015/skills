# Scenario

**Feature**: no arguments prints usage

```
user -> playwright-debug CLI (no args) -> usage
```

## Preconditions

- `req.Args` is empty.

## Steps

1. Set `req.Args = []string{}`.

## Context

- Exit code 0; usage on stdout.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{}
	return nil
}
```