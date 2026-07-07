# Scenario

**Feature**: `-h` prints usage

```
user -> playwright-debug CLI (-h) -> usage
```

## Preconditions

- Single help flag argument.

## Steps

1. Set `req.Args = []string{"-h"}`.

## Context

- Same content expectations as no-args help.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"-h"}
	return nil
}
```