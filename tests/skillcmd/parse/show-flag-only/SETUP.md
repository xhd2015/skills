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
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"--show"}
	return nil
}
```
