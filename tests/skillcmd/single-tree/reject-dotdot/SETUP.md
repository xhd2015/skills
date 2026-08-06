# Scenario

**Feature**: `--show ../x` rejects dotdot path segments

```
caller -> SingleSkill.Handle(--show ../x) -> error
```

## Preconditions

- Tree SingleSkill configured by parent.

## Steps

1. Configure Args with invalid `../x` topic path.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"--show", "../x"}
	return nil
}
```
