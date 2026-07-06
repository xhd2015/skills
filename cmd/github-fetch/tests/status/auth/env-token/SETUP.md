# Scenario

**Feature**: GITHUB_TOKEN takes priority over gh

```
# env token wins over gh CLI
GITHUB_TOKEN env + fake gh (logged in) -> auth resolver (GITHUB_TOKEN) -> mock API /user
```

## Preconditions

- `GITHUB_TOKEN` is set to `gho_secrettoken12345`.
- Fake logged-in `gh` is also on PATH (would authenticate as `ghuser` if used).
- Mock `/user` returns `testuser` when probed with env token.

## Steps

1. Set `req.GithubToken = "gho_secrettoken12345"`.
2. Set `req.GhMode = GhLoggedIn` with `req.GhUsername = "ghuser"` to prove token priority.
3. Run `github-fetch status`.

## Context

- Masked token display: `gho_****`.
- API access line: `authenticated (via GITHUB_TOKEN) as testuser`.

```go
func Setup(t *testing.T, req *Request) error {
	req.GithubToken = "gho_secrettoken12345"
	req.GhMode = GhLoggedIn
	req.GhUsername = "ghuser"
	req.GhHost = "github.com"
	req.GhScopes = "repo"
	req.GhToken = "gho_ghclitoken999"
	return nil
}
```