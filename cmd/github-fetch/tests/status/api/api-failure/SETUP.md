# Scenario

**Feature**: status fails when API probe returns 500

```
# API probe error is fatal
mock API 500 on /rate_limit or /user -> stderr error -> non-zero exit
```

## Preconditions

- Mock server returns HTTP 500 for all API endpoints.
- Unauthenticated setup (no gh, no token).

## Steps

1. Set `req.GhMode = GhAbsent`.
2. Set `req.MockAPIFail = true`.
3. Run `github-fetch status`.

## Context

- Informational auth lines may be partially printed before failure; primary assertion is stderr + exit code.

```go
func Setup(t *testing.T, req *Request) error {
	req.GhMode = GhAbsent
	req.MockAPIFail = true
	return nil
}
```