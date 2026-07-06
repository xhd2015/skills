# Scenario

**Feature**: `github-fetch status` reflects resolved credential source

```
# gh and/or GITHUB_TOKEN determine API auth method
gh auth status/token -> auth resolver -> status summary
GITHUB_TOKEN env -> auth resolver (priority over gh) -> status summary

# unauthenticated path still probes rate limits
auth resolver (none) -> mock API /rate_limit -> rate limit section
```

## Preconditions

- Mock API server returns default `/rate_limit` payload.
- `req.Args` defaults to `["status"]` from root setup.

## Steps

1. Each leaf sets `req.GhMode` and optional `req.GithubToken`.
2. Run `github-fetch status` and assert auth resolution lines in stdout.

## Context

- `no-gh-no-token` omits `gh` from PATH and leaves `GITHUB_TOKEN` unset.
- `env-token` sets both `GITHUB_TOKEN` and a logged-in fake `gh` to prove token priority.

```go
func Setup(t *testing.T, req *Request) error {
	if req.MockAPIFail {
		t.Fatalf("auth resolution leaves expect successful API probes")
	}
	if len(req.Args) == 0 || req.Args[0] != "status" {
		req.Args = []string{"status"}
	}
	return nil
}
```