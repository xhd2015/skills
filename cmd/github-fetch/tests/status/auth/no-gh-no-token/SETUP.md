# Scenario

**Feature**: status with no gh and no GITHUB_TOKEN

```
# no credential sources available
(no gh on PATH, no GITHUB_TOKEN) -> auth resolver (none) -> unauthenticated summary
mock API /rate_limit -> rate limit section
```

## Preconditions

- `gh` is not on PATH.
- `GITHUB_TOKEN` is not set.

## Steps

1. Set `req.GhMode = GhAbsent`.
2. Run `github-fetch status`.

## Context

- Informational output still exits 0 per requirement.

```go
func Setup(t *testing.T, req *Request) error {
	req.GhMode = GhAbsent
	return nil
}
```