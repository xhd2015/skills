# Scenario

**Feature**: `github-fetch status` reports auth source and rate limits via mock API

```
# status probes mock GitHub API with resolved credential
github-fetch status -> auth resolver -> mock API (/user, /rate_limit) -> structured summary

# gh CLI and GITHUB_TOKEN feed the resolver when present
gh auth status/token -> auth resolver
GITHUB_TOKEN env -> auth resolver (priority over gh)
```

## Preconditions

- `DOCTEST_ROOT` is `cmd/github-fetch/tests` (parent of this test tree).
- `findModuleRoot()` walks upward from `DOCTEST_ROOT` to locate the module `go.mod`.
- Tests use a mock HTTP server for `GET /user` and `GET /rate_limit`; no live GitHub calls.
- Fake `gh` scripts are prepended to `PATH` when `req.GhMode` is not `absent`.

## Steps

1. Build the `github-fetch` binary from `cmd/github-fetch` within the module.
2. Start the status mock HTTP server with `req.MockRateLimit` and failure mode.
3. Configure `PATH` and `GITHUB_TOKEN` per `req.GhMode`.
4. Set `GITHUB_API_BASE_URL` to the mock server URL.
5. Execute `github-fetch <req.Args...>`.
6. Capture stdout, stderr, exit code, and the mock server URL in `Response.APIBaseURL`.

## Context

- Default `req.Args` is `["status"]` when unset.
- Default mock rate limit: limit 5000, remaining 4987, reset 1751812200.
- Default mock user login: `testuser`.
- When `req.GhMode` is `absent`, `PATH` is set to an empty temp dir so `gh` is not found.

```go
import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func findModuleRoot(d *session.Doctest) (string, error) {
	dir := d.DOCTEST_ROOT
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("cannot find go.mod")
		}
		dir = parent
	}
}

func defaultMockRateLimit() MockRateLimit {
	return MockRateLimit{
		Limit:     5000,
		Remaining: 4987,
		Reset:     1751812200,
	}
}

func startStatusMockServer(req *Request) *httptest.Server {
	writeJSON := func(w http.ResponseWriter, status int, v interface{}) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(v)
	}

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if req.MockAPIFail {
			writeJSON(w, http.StatusInternalServerError, map[string]string{
				"message": "Internal Server Error",
			})
			return
		}

		switch {
		case r.Method == "GET" && r.URL.Path == "/user":
			writeJSON(w, http.StatusOK, map[string]string{
				"login": req.MockUserLogin,
			})
		case r.Method == "GET" && r.URL.Path == "/rate_limit":
			writeJSON(w, http.StatusOK, map[string]interface{}{
				"resources": map[string]interface{}{
					"core": map[string]interface{}{
						"limit":     req.MockRateLimit.Limit,
						"remaining": req.MockRateLimit.Remaining,
						"reset":     req.MockRateLimit.Reset,
					},
				},
			})
		default:
			writeJSON(w, http.StatusNotFound, map[string]string{
				"message": "Not Found",
			})
		}
	}))
}

func setupGhOnPath(t *testing.T, req *Request) (string, error) {
	if req.GhMode == "" {
		req.GhMode = GhAbsent
	}
	if req.GhMode == GhAbsent {
		return "", nil
	}

	binDir := filepath.Join(t.TempDir(), "bin")
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		return "", fmt.Errorf("create fake gh bin dir: %w", err)
	}

	ghScript, err := fakeGhScript(req)
	if err != nil {
		return "", err
	}

	ghPath := filepath.Join(binDir, "gh")
	if err := os.WriteFile(ghPath, []byte(ghScript), 0o755); err != nil {
		return "", fmt.Errorf("write fake gh: %w", err)
	}

	return binDir, nil
}

func fakeGhScript(req *Request) (string, error) {
	host := req.GhHost
	if host == "" {
		host = "github.com"
	}
	username := req.GhUsername
	if username == "" {
		username = "testuser"
	}
	scopes := req.GhScopes
	if scopes == "" {
		scopes = "repo, read:org"
	}
	token := req.GhToken
	if token == "" {
		token = "gho_fakeghtoken123"
	}

	switch req.GhMode {
	case GhLoggedIn:
		return fmt.Sprintf(`#!/bin/sh
case "$1" in
auth)
  case "$2" in
  status)
    echo "%s"
    echo "Logged in to %s as %s"
    echo "Token scopes: %s"
    exit 0
    ;;
  token)
    echo "%s"
    exit 0
    ;;
  esac
  ;;
esac
echo "unknown gh command: $*" >&2
exit 1
`, host, host, username, scopes, token), nil
	case GhNotLoggedIn:
		return `#!/bin/sh
case "$1" in
auth)
  case "$2" in
  status)
    echo "not logged in" >&2
    exit 1
    ;;
  token)
    echo "not logged in" >&2
    exit 1
    ;;
  esac
  ;;
esac
echo "unknown gh command: $*" >&2
exit 1
`, nil
	default:
		return "", fmt.Errorf("unsupported gh mode: %q", req.GhMode)
	}
}

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	if len(req.Args) == 0 {
		req.Args = []string{"status"}
	}
	if req.MockUserLogin == "" {
		req.MockUserLogin = "testuser"
	}
	if req.MockRateLimit.Limit == 0 && req.MockRateLimit.Remaining == 0 && req.MockRateLimit.Reset == 0 {
		req.MockRateLimit = defaultMockRateLimit()
	}
	return nil
}
```