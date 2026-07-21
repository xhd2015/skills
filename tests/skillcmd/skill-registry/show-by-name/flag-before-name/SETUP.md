# Scenario

**Feature**: `--show foo` prints foo skill content

```
# flag before name
caller -> HandleSkill(--show foo) -> Foo Skill Body
```

## Steps

1. Set Args `--show foo`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"--show", "foo"}
	return nil
}
```
