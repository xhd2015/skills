# Scenario

**Feature**: status when gh is available but not logged in

```
# gh present but unauthenticated
fake gh (not logged in) -> auth resolver (none) -> unauthenticated summary
mock API /rate_limit -> rate limit section
```

## Preconditions

- Fake `gh` on PATH exits non-zero for `gh auth status`.
- `GITHUB_TOKEN` is not set.

## Steps

1. Set `req.GhMode = GhNotLoggedIn`.
2. Run `github-fetch status`.

## Context

- gh CLI line should read `available`; API access remains unauthenticated.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.GhMode = GhNotLoggedIn
	return nil
}
```