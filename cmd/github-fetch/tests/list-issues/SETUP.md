## Preconditions
- `DOCTEST_ROOT` is `cmd/github-fetch/tests` (the parent of this test tree).
- `findModuleRoot()` walks upward from `DOCTEST_ROOT` to locate the module `go.mod`.
- A mock GitHub API server replaces the real GitHub API for deterministic tests.

## Steps
1. Build the `github-fetch` binary from `cmd/github-fetch` within the module.
2. If `req.InGitRepo` is true: create a temporary git repository with a single commit, configure the specified origin remote URL, and `cd` into it.
3. Start a mock HTTP server configured with `req.MockIssues`.
4. Set `GITHUB_API_BASE_URL` to the mock server URL before invoking the binary.
5. Execute `github-fetch <req.Args...>`.
6. Capture stdout, stderr, and exit code.

## Context
- The `github-fetch` binary is built from `cmd/github-fetch` within this module.
- The mock server handles `GET /repos/{owner}/{repo}` and `GET /repos/{owner}/{repo}/issues` with state filtering and pagination.
- Items with a `pull_request` field in the JSON response are excluded from issue list output (GitHub returns PRs in the issues endpoint).

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

type MockIssueListItem struct {
	Number        int
	Title         string
	State         string
	User          string
	HTMLURL       string
	Labels        []string
	IsPullRequest bool
}

type Request struct {
	Args        []string
	OriginURL   string
	InGitRepo   bool
	MockIssues  []MockIssueListItem
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

func startIssueMockServer(issues []MockIssueListItem) *httptest.Server {
	writeJSON := func(w http.ResponseWriter, v interface{}) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(v)
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

		case len(parts) == 4 && parts[3] == "issues" && r.Method == "GET":
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

			var filtered []MockIssueListItem
			for _, issue := range issues {
				if state == "all" || issue.State == state {
					filtered = append(filtered, issue)
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

			type issueLabel struct {
				Name string `json:"name"`
			}
			type issueUser struct {
				Login string `json:"login"`
			}
			type pullRequestStub struct {
				URL string `json:"url"`
			}
			type issueResp struct {
				Number      int               `json:"number"`
				Title       string            `json:"title"`
				State       string            `json:"state"`
				HTMLURL     string            `json:"html_url"`
				User        issueUser         `json:"user"`
				Labels      []issueLabel      `json:"labels"`
				PullRequest *pullRequestStub  `json:"pull_request,omitempty"`
			}

			paged := make([]issueResp, len(pageItems))
			for i, item := range pageItems {
				resp := issueResp{
					Number: item.Number, Title: item.Title, State: item.State,
					HTMLURL: item.HTMLURL, User: issueUser{Login: item.User},
				}
				for _, label := range item.Labels {
					resp.Labels = append(resp.Labels, issueLabel{Name: label})
				}
				if item.IsPullRequest {
					resp.PullRequest = &pullRequestStub{URL: item.HTMLURL}
				}
				paged[i] = resp
			}

			writeJSON(w, paged)

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

	mockSrv := startIssueMockServer(req.MockIssues)
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