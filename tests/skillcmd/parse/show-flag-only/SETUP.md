# Scenario

**Feature**: `--show` alone selects show with empty rest

```
caller -> ParseSkillArgs([--show]) -> Action=show, Rest=[]
```

## Preconditions

- Parse mode configured by parent.

## Steps

1. Set Args for this case.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{"--show"}
	return nil
}
```
