# Scenario

**Feature**: status API base URL display and probe error handling

```
# custom API host from env
GITHUB_API_BASE_URL -> status summary API base URL line

# probe failures surface as command errors
mock API 500 -> stderr error, non-zero exit
```

## Preconditions

- `GITHUB_API_BASE_URL` is set to the mock server URL by root `Run`.
- Leaves configure mock success vs failure behavior.

## Steps

1. `custom-api-base` uses default mock responses and asserts URL echo.
2. `api-failure` sets `req.MockAPIFail = true`.

## Context

- `custom-api-base` uses unauthenticated path (no gh, no token) for simplicity.

```go
func Setup(t *testing.T, req *Request) error {
	if len(req.Args) == 0 || req.Args[0] != "status" {
		req.Args = []string{"status"}
	}
	if req.GhMode == "" {
		req.GhMode = GhAbsent
	}
	return nil
}
```