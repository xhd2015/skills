# Scenario

**Feature**: top-level `install` remains backward-compatible with `skill --install`

```
# legacy entry point still routes to install handler
user -> go-best-practice install --dry-run -> [dry-run] .agents/skills/go-best-practice
```

## Preconditions

- Command runs in an isolated work directory.

## Steps

1. Set `req.UseWorkDir = true`.
2. Set `req.Args = ["install", "--dry-run"]` (no `skill` prefix).

```go
func Setup(t *testing.T, req *Request) error {
	req.UseWorkDir = true
	req.Args = []string{"install", "--dry-run"}
	return nil
}
```
