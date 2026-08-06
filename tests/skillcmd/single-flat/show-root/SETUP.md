# Scenario

**Feature**: `--show` prints RootContent

```
caller -> SingleSkill.Handle(--show) -> RootContent on stdout
```

## Preconditions

- Flat SingleSkill configured by parent.

## Steps

1. Set Args for this action.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"--show"}
	return nil
}
```
