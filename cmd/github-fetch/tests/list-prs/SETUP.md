## Preconditions
- `DOCTEST_ROOT` is `cmd/github-fetch/tests` (the parent of this test tree).
- `findModuleRoot()` walks upward from `DOCTEST_ROOT` to locate the module `go.mod`.
- A mock GitHub API server replaces the real GitHub API for deterministic tests.

## Steps
1. Build the `github-fetch` binary from `cmd/github-fetch` within the module.
2. If `req.InGitRepo` is true: create a temporary git repository with a single commit, configure the specified origin remote URL, and `cd` into it.
3. Start a mock HTTP server configured with `req.MockPRs` and `req.HasNextPage`.
4. Set `GITHUB_API_BASE_URL` to the mock server URL before invoking the binary.
5. Execute `github-fetch <req.Args...>`.
6. Capture stdout, stderr, and exit code.

## Context
- The `github-fetch` binary is built from `cmd/github-fetch` within this module.
- The mock server handles `GET /repos/{owner}/{repo}` and `GET /repos/{owner}/{repo}/pulls` with state filtering, pagination, and optional `Link: rel="next"` header.
- When `req.HasNextPage` is true and page 1 is requested, the mock sets a `Link` header with `rel="next"`.
- Default `--per-page` is 30; the mock honors `page` and `per_page` query parameters.

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
	"strconv"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/session"
)

type MockPRListItem struct {
	Number  int
	Title   string
	State   string
	User    string
	HTMLURL string
}

type Request struct {
	Args        []string
	OriginURL   string
	InGitRepo   bool
	MockPRs     []MockPRListItem
	HasNextPage bool
}

type Response struct {
	ExitCode int
	Stdout   string
	Stderr   string
}

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

func startPRMockServer(prs []MockPRListItem, hasNextPage bool) *httptest.Server {
	writeJSON := func(w http.ResponseWriter, v interface{}) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(v)
	}
	scheme := func(r *http.Request) string {
		if r.TLS != nil {
			return "https"
		}
		return "http"
	}

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/")
		parts := strings.Split(path, "/")

		switch {
		case len(parts) == 3 && r.Method == "GET":
			owner, repo := parts[1], parts[2]
			writeJSON(w, map[string]string{
				"full_name":      owner + "/" + repo,
				"default_branch": "main",
			})

		case len(parts) == 4 && parts[3] == "pulls" && r.Method == "GET":
			state := r.URL.Query().Get("state")
			if state == "" {
				state = "open"
			}
			page, _ := strconv.Atoi(r.URL.Query().Get("page"))
			if page < 1 {
				page = 1
			}
			perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
			if perPage < 1 {
				perPage = 30
			}

			var filtered []MockPRListItem
			for _, pr := range prs {
				if state == "all" || pr.State == state {
					filtered = append(filtered, pr)
				}
			}

			start := (page - 1) * perPage
			if start > len(filtered) {
				start = len(filtered)
			}
			end := start + perPage
			if end > len(filtered) {
				end = len(filtered)
			}
			pageItems := filtered[start:end]

			type prUser struct {
				Login string `json:"login"`
			}
			type prResp struct {
				Number  int    `json:"number"`
				Title   string `json:"title"`
				State   string `json:"state"`
				HTMLURL string `json:"html_url"`
				User    prUser `json:"user"`
			}
			out := make([]prResp, len(pageItems))
			for i, pr := range pageItems {
				out[i] = prResp{
					Number: pr.Number, Title: pr.Title, State: pr.State,
					HTMLURL: pr.HTMLURL, User: prUser{Login: pr.User},
				}
			}

			if hasNextPage && page == 1 && end < len(filtered) {
				nextURL := fmt.Sprintf("%s://%s%s?page=%d&per_page=%d&state=%s",
					scheme(r), r.Host, r.URL.Path, page+1, perPage, state)
				w.Header().Set("Link", fmt.Sprintf(`<%s>; rel="next"`, nextURL))
			}

			writeJSON(w, out)

		default:
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"message": "Not Found"})
		}
	}))
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	moduleRoot, err := findModuleRoot(d)
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

	var runDir string
	if req.InGitRepo {
		runDir, err = os.MkdirTemp("", "github-fetch-test-*")
		if err != nil {
			return nil, fmt.Errorf("create temp dir: %w", err)
		}
		defer os.RemoveAll(runDir)

		if out, err := runInDir(runDir, "git", "init"); err != nil {
			return nil, fmt.Errorf("git init: %v\n%s", err, out)
		}
		if err := os.WriteFile(filepath.Join(runDir, "README.md"), []byte("# test"), 0644); err != nil {
			return nil, fmt.Errorf("write file: %w", err)
		}
		runInDir(runDir, "git", "config", "user.email", "test@test.com")
		runInDir(runDir, "git", "config", "user.name", "test")
		if out, err := runInDir(runDir, "git", "add", "."); err != nil {
			return nil, fmt.Errorf("git add: %v\n%s", err, out)
		}
		if out, err := runInDir(runDir, "git", "commit", "-m", "initial"); err != nil {
			return nil, fmt.Errorf("git commit: %v\n%s", err, out)
		}

		if req.OriginURL != "" {
			if out, err := runInDir(runDir, "git", "remote", "add", "origin", req.OriginURL); err != nil {
				return nil, fmt.Errorf("git remote add origin: %v\n%s", err, out)
			}
		}
	} else {
		runDir, err = os.MkdirTemp("", "github-fetch-test-*")
		if err != nil {
			return nil, fmt.Errorf("create temp dir: %w", err)
		}
		defer os.RemoveAll(runDir)
	}

	mockSrv := startPRMockServer(req.MockPRs, req.HasNextPage)
	defer mockSrv.Close()

	cmd := exec.Command(binaryPath, req.Args...)
	cmd.Dir = runDir
	cmd.Env = append(os.Environ(), "GITHUB_API_BASE_URL="+mockSrv.URL)

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
		ExitCode: exitCode,
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
	}, nil
}

func runInDir(dir string, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	return string(out), err
}
```