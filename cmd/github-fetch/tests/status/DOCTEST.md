# Status Command Tests

Doc-style tests for `github-fetch status`: reports GitHub authentication status,
resolved API credential source, and core rate-limit information from the GitHub API.

# DSN (Domain Specific Notion)

Participants:

- The `github-fetch` CLI binary exposes a `status` sub-command that prints a
  structured auth and rate-limit summary to stdout.
- The `gh` CLI (when on PATH) supplies login state via `gh auth status` and a
  token via `gh auth token`.
- The `GITHUB_TOKEN` environment variable is a first-class auth source with
  higher priority than `gh` for all API probes.
- A mock GitHub REST API server stands in for `api.github.com`, serving
  `GET /user` and `GET /rate_limit` without calling the live GitHub API.
- Auth resolution picks one of: `GITHUB_TOKEN`, `gh`, or unauthenticated HTTP;
  `status` probes the API with that credential and prints which method is active.

## Decision Tree

```text
cmd/github-fetch/tests/status/
├── [auth resolution]
│   └── auth/
│       ├── no-gh-no-token/       # PATH without gh, no GITHUB_TOKEN
│       ├── gh-logged-in/         # fake gh logged in → via gh
│       ├── gh-not-logged-in/     # fake gh present but not logged in
│       └── env-token/            # GITHUB_TOKEN beats gh
├── [API interaction]
│   └── api/
│       ├── custom-api-base/      # GITHUB_API_BASE_URL echoed in output
│       └── api-failure/          # mock 500 → stderr error, non-zero exit
└── [help]
    └── help/
        ├── status-help/            # status -h usage
        └── root-help/              # github-fetch -h lists status
```

## Test Index

| # | Leaf | Description |
|---|------|-------------|
| 1 | auth/no-gh-no-token | No gh, no token → unauthenticated summary + rate limit |
| 2 | auth/gh-logged-in | Fake gh logged in → authenticated via gh |
| 3 | auth/gh-not-logged-in | Fake gh not logged in → unauthenticated |
| 4 | auth/env-token | GITHUB_TOKEN set → authenticated via GITHUB_TOKEN (over gh) |
| 5 | api/custom-api-base | Custom `GITHUB_API_BASE_URL` shown in output |
| 6 | api/api-failure | API probe failure → stderr error, non-zero exit |
| 7 | help/status-help | `status -h` prints status usage |
| 8 | help/root-help | Root `-h` lists `status` command |

## How to Run

```sh
doctest vet ./cmd/github-fetch/tests/status
doctest test -v ./cmd/github-fetch/tests/status
```

## Version

0.0.2

```go
import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

type GhMode string

const (
	GhAbsent      GhMode = "absent"
	GhLoggedIn    GhMode = "logged-in"
	GhNotLoggedIn GhMode = "not-logged-in"
)

type MockRateLimit struct {
	Limit     int
	Remaining int
	Reset     int64
}

type Request struct {
	Args []string

	GhMode     GhMode
	GhHost     string
	GhUsername string
	GhScopes   string
	GhToken    string

	GithubToken string

	MockUserLogin string
	MockRateLimit MockRateLimit
	MockAPIFail   bool
}

type Response struct {
	ExitCode   int
	Stdout     string
	Stderr     string
	APIBaseURL string
}

func Run(t *testing.T, req *Request) (*Response, error) {
	moduleRoot, err := findModuleRoot()
	if err != nil {
		return nil, fmt.Errorf("find module root: %w", err)
	}

	binaryPath := filepath.Join(t.TempDir(), "github-fetch")
	buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
	buildCmd.Dir = filepath.Join(moduleRoot, "cmd", "github-fetch")
	buildOut, err := buildCmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("build github-fetch: %v\n%s", err, string(buildOut))
	}

	mockSrv := startStatusMockServer(req)
	defer mockSrv.Close()

	runDir, err := os.MkdirTemp("", "github-fetch-status-*")
	if err != nil {
		return nil, fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(runDir)

	env := os.Environ()
	env = append(env, "GITHUB_API_BASE_URL="+mockSrv.URL)
	if req.GithubToken != "" {
		env = append(env, "GITHUB_TOKEN="+req.GithubToken)
	}

	pathPrefix, err := setupGhOnPath(t, req)
	if err != nil {
		return nil, err
	}
	if pathPrefix != "" {
		env = append(env, "PATH="+pathPrefix+string(os.PathListSeparator)+os.Getenv("PATH"))
	} else {
		env = append(env, "PATH="+t.TempDir())
	}

	cmd := exec.Command(binaryPath, req.Args...)
	cmd.Dir = runDir
	cmd.Env = env

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	runErr := cmd.Run()
	exitCode := 0
	if runErr != nil {
		if exitErr, ok := runErr.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			return nil, fmt.Errorf("run github-fetch: %w", runErr)
		}
	}

	return &Response{
		ExitCode:   exitCode,
		Stdout:     stdout.String(),
		Stderr:     stderr.String(),
		APIBaseURL: mockSrv.URL,
	}, nil
}
```