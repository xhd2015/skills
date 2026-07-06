# Scenario

**Feature**: status echoes custom GITHUB_API_BASE_URL

```
# env overrides default api.github.com
GITHUB_API_BASE_URL=mock URL -> status summary shows API base URL line
```

## Preconditions

- `GITHUB_API_BASE_URL` points at the mock httptest server.
- No gh and no GITHUB_TOKEN (informational unauthenticated output).

## Steps

1. Set `req.GhMode = GhAbsent`.
2. Run `github-fetch status`.
3. Assert `API base URL` line contains the mock server URL.

## Context

- Verifies the binary reads and displays `GITHUB_API_BASE_URL`.

```go
func Setup(t *testing.T, req *Request) error {
	req.GhMode = GhAbsent
	return nil
}
```