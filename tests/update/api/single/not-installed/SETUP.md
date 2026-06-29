# Scenario

**Feature**: update when no target has `SKILL.md`

```
HandleUpdate -> each resolved dir missing SKILL.md -> silent skip
```

## Preconditions

- Leaves do not run `PreInstalls`.

## Steps

1. Configure `SingleOpts` and call update.

## Context

- Sibling leaves cover the zero-install case only.

```go
func Setup(t *testing.T, req *Request) error {
	req.PreInstalls = nil
	return nil
}
```