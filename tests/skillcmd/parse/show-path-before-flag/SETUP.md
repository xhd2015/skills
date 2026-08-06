# Scenario

**Feature**: path before `--show` lands in Rest

```
caller -> ParseSkillArgs([flags-parsing/types, --show]) -> show + rest path
```

## Preconditions

- Parse mode configured by parent.

## Steps

1. Set Args for this case.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"flags-parsing/types", "--show"}
	return nil
}
```
