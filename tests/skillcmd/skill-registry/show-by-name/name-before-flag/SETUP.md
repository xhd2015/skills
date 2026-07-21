# Scenario

**Feature**: `foo --show` prints the same foo skill content

```
# name before flag (both orders)
caller -> HandleSkill(foo --show) -> Foo Skill Body
```

## Steps

1. Set Args `foo --show`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"foo", "--show"}
	return nil
}
```
