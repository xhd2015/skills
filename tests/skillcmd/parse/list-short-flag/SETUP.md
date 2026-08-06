# Scenario

**Feature**: `-l` is an alias of `--list`

```
caller -> ParseSkillArgs([-l]) -> Action=list
```

## Preconditions

- Parse mode configured by parent.

## Steps

1. Set Args for this case.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"-l"}
	return nil
}
```
