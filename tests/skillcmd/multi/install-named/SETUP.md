# Scenario

**Feature**: `--install foo --dry-run` targets the named skill directory

```
# install by name with dry-run
caller -> HandleSkill(--install foo --dry-run) -> .agents/skills/foo
```

## Steps

1. Use workdir isolation.
2. Set Args for named install dry-run.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.UseWorkDir = true
	req.Args = []string{"--install", "foo", "--dry-run"}
	return nil
}
```
