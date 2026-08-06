# Scenario

**Feature**: `--show --header` sets Header true

```
caller -> ParseSkillArgs([--show, --header]) -> show + Header
```

## Preconditions

- Parse mode configured by parent.

## Steps

1. Set Args for this case.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"--show", "--header"}
	return nil
}
```
