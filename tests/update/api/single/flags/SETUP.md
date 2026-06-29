# Scenario

**Feature**: update command help documents supported flags

```
HandleUpdate --help -> usage text on stdout
```

## Preconditions

- No pre-install required.

## Steps

1. Leaves pass `-h` or `--help` in `UpdateArgs`.

## Context

- Help must not mutate filesystem.

```go
func Setup(t *testing.T, req *Request) error {
	req.PreInstalls = nil
	return nil
}
```