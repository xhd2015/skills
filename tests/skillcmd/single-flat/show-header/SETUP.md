# Scenario

**Feature**: `--show --header` prints frontmatter only

```
caller -> SingleSkill.Handle(--show --header) -> FormatHeaderWithDelimiters
```

## Preconditions

- Flat SingleSkill configured by parent.

## Steps

1. Set Args for this action.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"--show", "--header"}
	return nil
}
```
