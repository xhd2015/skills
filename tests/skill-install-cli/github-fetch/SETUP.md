# Scenario

**Feature**: `github-fetch skill --install` global dry-run and install routing rules

```
# skill --install resolves global targets under HOME
user -> github-fetch skill --install --global --dry-run -> HOME/.agents/skills/github-fetch

# top-level install is not a valid command
user -> github-fetch install -> unknown command error
```

## Preconditions

- `github-fetch` binary is built from `cmd/github-fetch`.

## Steps

1. Resolve `req.Binary` via session cache build.
2. Set `req.SkillName = "github-fetch"`.
3. Leaves configure args for success or error paths.

## Context

- Unlike `go-best-practice`, `github-fetch` has no top-level `install` alias.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	bin, err := buildGithubFetchOnce(t, d)
	if err != nil {
		return err
	}
	req.Binary = bin
	req.SkillName = "github-fetch"
	return nil
}
```
