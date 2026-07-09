# Scenario

**Feature**: combining `--show` and `--install` is an error

```
caller -> ParseSkillArgs([--show, --install]) -> error
```

## Preconditions

- Parse mode configured by parent.

## Steps

1. Set Args for this case.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"--show", "--install"}
	return nil
}
```
