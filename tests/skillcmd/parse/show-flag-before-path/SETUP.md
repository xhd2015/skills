# Scenario

**Feature**: `--show` before path lands path in Rest

```
caller -> ParseSkillArgs([--show, flags-parsing/types]) -> show + rest path
```

## Preconditions

- Parse mode configured by parent.

## Steps

1. Set Args for this case.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"--show", "flags-parsing/types"}
	return nil
}
```
