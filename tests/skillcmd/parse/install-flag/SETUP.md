# Scenario

**Feature**: `--install --global` keeps install options in Rest

```
caller -> ParseSkillArgs([--install, --global]) -> install + rest --global
```

## Preconditions

- Parse mode configured by parent.

## Steps

1. Set Args for this case.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"--install", "--global"}
	return nil
}
```
