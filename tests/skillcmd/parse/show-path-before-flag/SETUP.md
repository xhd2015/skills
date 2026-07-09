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
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"flags-parsing/types", "--show"}
	return nil
}
```
