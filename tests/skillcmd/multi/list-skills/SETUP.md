# Scenario

**Feature**: `--list` prints registered skill names and descriptions

```
# list registry entries in registration order
caller -> HandleSkill(--list) -> names (+ descriptions)
```

## Preconditions

- Registry has foo and bar.

## Steps

1. Set Args to `--list`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"--list"}
	return nil
}
```
