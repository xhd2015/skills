# Scenario

**Feature**: global dry-run resolves install target under HOME

```
# HOME is isolated temp dir for global scope
user -> go-best-practice skill --install --global --dry-run -> HOME/.agents/skills/go-best-practice
```

## Preconditions

- `HOME` is set to a temporary directory.

## Steps

1. Set `req.UseGlobalHome = true`.
2. Set `req.Args = ["skill", "--install", "--global", "--dry-run"]`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.UseGlobalHome = true
	req.Args = []string{"skill", "--install", "--global", "--dry-run"}
	return nil
}
```
