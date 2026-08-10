# Scenario

**Feature**: status when fake gh is logged in

```
# gh provides token for API probes
fake gh (logged in) -> auth resolver (gh) -> mock API /user + /rate_limit -> authenticated via gh
```

## Preconditions

- Fake `gh` on PATH reports logged-in state for `testuser`.
- Mock `/user` returns `{"login": "testuser"}`.

## Steps

1. Set `req.GhMode = GhLoggedIn`.
2. Run `github-fetch status`.

## Context

- Expected API access line: `authenticated (via gh) as testuser`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.GhMode = GhLoggedIn
	req.GhUsername = "testuser"
	req.GhHost = "github.com"
	req.GhScopes = "repo, read:org"
	return nil
}
```