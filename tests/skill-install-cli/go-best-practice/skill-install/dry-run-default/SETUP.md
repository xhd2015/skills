# Scenario

**Feature**: local dry-run previews default `.agents` target

```
# default skill install dry-run in isolated work dir
user -> go-best-practice skill install --dry-run -> [dry-run] .agents/skills/go-best-practice
```

## Preconditions

- Command runs in an isolated work directory.

## Steps

1. Set `req.UseWorkDir = true`.
2. Set `req.Args = ["skill", "install", "--dry-run"]`.

```go
func Setup(t *testing.T, req *Request) error {
	req.UseWorkDir = true
	req.Args = []string{"skill", "install", "--dry-run"}
	return nil
}
```