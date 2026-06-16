## Preconditions
- `DOCTEST_ROOT` is `cmd/github-fetch/tests` (the parent of this test tree).
- `findModuleRoot()` walks upward from `DOCTEST_ROOT` to locate the `github.com/xhd2015/skills` module `go.mod`.
- A mock GitHub API server replaces the real GitHub API for deterministic tests.

## Steps
1. Build the `github-fetch` binary from `cmd/github-fetch` within the module.
2. If `req.InGitRepo` is true: create a temporary git repository with a single commit, configure the specified origin remote URL, and `cd` into it.
3. Start a mock HTTP server using the `githubmock` package, configuring it with `req.MockRuns`, `req.MockJobs`, `req.MockLogs`, and `req.MockWorkflowFiles`.
4. Set the environment variable `GITHUB_API_BASE_URL` to the mock server URL before invoking the binary.
5. Execute `github-fetch ci <req.Args...>`.
6. Capture stdout, stderr, and exit code.

## Context
- The `github-fetch` binary is built from `cmd/github-fetch` within this module.
- The mock server URL is passed via `GITHUB_API_BASE_URL` environment variable.
- All GitHub API calls within the binary must use this base URL instead of the hardcoded `https://api.github.com`.
- The default branch is returned by the mock via `GET /repos/{owner}/{repo}` endpoint.

```go
import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/xhd2015/skills/cmd/github-fetch/github-mock"
)

type MockWorkflowRun struct {
	ID         int64
	Name       string
	Status     string
	Conclusion string
	HTMLURL    string
	HeadBranch string
	Event      string
}

type MockWorkflowJob struct {
	ID         int64
	Name       string
	Conclusion string
}

type Request struct {
	Args      []string
	OriginURL string // git remote origin URL, e.g. "git@github.com:owner/repo.git"
	InGitRepo bool   // create a git repo and cd into it before running

	// Mock API response configuration
	MockRuns          []MockWorkflowRun
	MockJobs          []MockWorkflowJob
	MockLogs          string
	MockWorkflowFiles []string
	MockDefaultBranch string // override default branch (default: "main")
}

type Response struct {
	ExitCode int
	Stdout   string
	Stderr   string
}

func findModuleRoot() (string, error) {
	dir := DOCTEST_ROOT
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

	mockSrv := githubmock.NewServer()
	defer mockSrv.Close()

	mockSrv.WorkflowRuns = make([]githubmock.WorkflowRun, len(req.MockRuns))
	for i, r := range req.MockRuns {
		mockSrv.WorkflowRuns[i] = githubmock.WorkflowRun{
			ID:         r.ID,
			Name:       r.Name,
			Status:     r.Status,
			Conclusion: r.Conclusion,
			HTMLURL:    r.HTMLURL,
			HeadBranch: r.HeadBranch,
			Event:      r.Event,
		}
	}
	mockSrv.WorkflowJobs = make([]githubmock.WorkflowJob, len(req.MockJobs))
	for i, j := range req.MockJobs {
		mockSrv.WorkflowJobs[i] = githubmock.WorkflowJob{
			ID:         j.ID,
			Name:       j.Name,
			Conclusion: j.Conclusion,
		}
	}
	mockSrv.JobLogs = req.MockLogs
	mockSrv.WorkflowFiles = req.MockWorkflowFiles
	if req.MockDefaultBranch != "" {
		mockSrv.DefaultBranch = req.MockDefaultBranch
	}

	cmdArgs := append([]string{"ci"}, req.Args...)
	cmd := exec.Command(binaryPath, cmdArgs...)
	cmd.Dir = runDir
	cmd.Env = append(os.Environ(), "GITHUB_API_BASE_URL="+mockSrv.URL())

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