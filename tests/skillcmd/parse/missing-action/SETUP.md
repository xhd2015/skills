# Scenario

**Feature**: bare path without action flag is an error

```
caller -> ParseSkillArgs([foo]) -> error (missing action)
```

## Preconditions

- Parse mode configured by parent.

## Steps

1. Set Args for this case.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"foo"}
	return nil
}
```
