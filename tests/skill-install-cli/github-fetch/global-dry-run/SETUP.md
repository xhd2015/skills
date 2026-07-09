# Scenario

**Feature**: github-fetch global dry-run resolves under HOME

```
# HOME is isolated temp dir for global scope
user -> github-fetch skill --install --global --dry-run -> HOME/.agents/skills/github-fetch
```

## Preconditions

- `HOME` is set to a temporary directory.

## Steps

1. Set `req.UseGlobalHome = true`.
2. Set `req.Args = ["skill", "--install", "--global", "--dry-run"]`.

```go
func Setup(t *testing.T, req *Request) error {
	req.UseGlobalHome = true
	req.Args = []string{"skill", "--install", "--global", "--dry-run"}
	return nil
}
```
