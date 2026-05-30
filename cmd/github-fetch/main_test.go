package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseGitHubURL(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		owner   string
		repo    string
		number  string
		wantErr bool
	}{
		{"full url", "https://github.com/xhd2015/xgo/pull/379", "xhd2015", "xgo", "379", false},
		{"with trailing slash", "https://github.com/owner/repo/pull/123/", "owner", "repo", "123", false},
		{"with files path", "https://github.com/owner/repo/pull/123/files", "owner", "repo", "123", false},
		{"http scheme", "http://github.com/owner/repo/pull/456", "owner", "repo", "456", false},
		{"not a github url", "https://gitlab.com/owner/repo/pull/123", "", "", "", true},
		{"not a PR url", "https://github.com/owner/repo/issues/123", "", "", "", true},
		{"missing number", "https://github.com/owner/repo/pull/", "", "", "", true},
		{"random string", "not-a-url", "", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			owner, repo, number, err := parseGitHubURL(tt.raw)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error for %q, got nil", tt.raw)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if owner != tt.owner {
				t.Errorf("owner = %q, want %q", owner, tt.owner)
			}
			if repo != tt.repo {
				t.Errorf("repo = %q, want %q", repo, tt.repo)
			}
			if number != tt.number {
				t.Errorf("number = %q, want %q", number, tt.number)
			}
		})
	}
}

func TestParseOriginURL(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		owner   string
		repo    string
		wantErr bool
	}{
		{"ssh url", "git@github.com:xhd2015/xgo.git", "xhd2015", "xgo", false},
		{"ssh url without .git", "git@github.com:xhd2015/xgo", "xhd2015", "xgo", false},
		{"https url", "https://github.com/xhd2015/xgo.git", "xhd2015", "xgo", false},
		{"https url without .git", "https://github.com/xhd2015/xgo", "xhd2015", "xgo", false},
		{"http url", "http://github.com/owner/repo", "owner", "repo", false},
		{"ssh:// url", "ssh://git@github.com/owner/repo.git", "owner", "repo", false},
		{"ssh:// url without .git", "ssh://git@github.com/owner/repo", "owner", "repo", false},
		{"invalid url", "not-a-url", "", "", true},
		{"gitlab url", "git@gitlab.com:owner/repo.git", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			owner, repo, err := parseOriginURL(tt.raw)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error for %q, got nil", tt.raw)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if owner != tt.owner {
				t.Errorf("owner = %q, want %q", owner, tt.owner)
			}
			if repo != tt.repo {
				t.Errorf("repo = %q, want %q", repo, tt.repo)
			}
		})
	}
}

func TestParsePRBranch(t *testing.T) {
	tests := []struct {
		branch string
		number string
		ok     bool
	}{
		{"pr-379", "379", true},
		{"pr-1", "1", true},
		{"pr-12345", "12345", true},
		{"main", "", false},
		{"feature/pr-379", "", false},
		{"pr-abc", "", false},
		{"PR-379", "", false},
		{"", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.branch, func(t *testing.T) {
			num, ok := parsePRBranch(tt.branch)
			if ok != tt.ok {
				t.Errorf("ok = %v, want %v", ok, tt.ok)
			}
			if num != tt.number {
				t.Errorf("number = %q, want %q", num, tt.number)
			}
		})
	}
}

func TestResolvePRRefWithNumber(t *testing.T) {
	dir := initGitRepo(t, "https://github.com/testowner/testrepo.git")
	oldDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(oldDir)

	owner, repo, number, err := resolvePRRef("42")
	if err != nil {
		t.Fatalf("resolvePRRef(42): %v", err)
	}
	if owner != "testowner" {
		t.Errorf("owner = %q, want testowner", owner)
	}
	if repo != "testrepo" {
		t.Errorf("repo = %q, want testrepo", repo)
	}
	if number != "42" {
		t.Errorf("number = %q, want 42", number)
	}
}

func TestVerifyOriginMatches(t *testing.T) {
	dir := initGitRepo(t, "git@github.com:testowner/testrepo.git")
	oldDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(oldDir)

	err := verifyOriginMatches(PRRepo{FullName: "testowner/testrepo"})
	if err != nil {
		t.Errorf("expected match, got error: %v", err)
	}

	err = verifyOriginMatches(PRRepo{FullName: "other/repo"})
	if err == nil {
		t.Error("expected mismatch error, got nil")
	}
}

func TestGetCurrentBranch(t *testing.T) {
	dir := initGitRepo(t, "https://github.com/owner/repo.git")
	oldDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(oldDir)

	// Default branch after init
	branch, err := getCurrentBranch()
	if err != nil {
		t.Fatalf("getCurrentBranch: %v", err)
	}
	if branch == "" {
		t.Error("expected non-empty branch name")
	}
}

func TestFetchPRInfoWithMock(t *testing.T) {
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/files") {
			w.Write([]byte(`[{"filename":"main.go","status":"modified","additions":5,"deletions":3,"changes":8,"patch":"@@ -1 +1 @@\n-old\n+new"}]`))
			return
		}
		resp := githubPRResponse{
			Number:  379,
			Title:   "Test PR",
			Body:    "This is a test",
			State:   "open",
			HTMLURL: "https://github.com/testowner/testrepo/pull/379",
			Head: githubPRBranch{
				Ref: "feature-branch",
				SHA: "abc123",
				Repo: githubPRRepo{
					FullName: "forkowner/testrepo",
					CloneURL: "https://github.com/forkowner/testrepo.git",
					SSHURL:   "git@github.com:forkowner/testrepo.git",
				},
			},
			Base: githubPRBranch{
				Ref: "main",
				SHA: "def456",
				Repo: githubPRRepo{
					FullName: "testowner/testrepo",
					CloneURL: "https://github.com/testowner/testrepo.git",
					SSHURL:   "git@github.com:testowner/testrepo.git",
				},
			},
			User: githubPRUser{Login: "contributor"},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer mock.Close()

	// Override the API base URL - we can't, so we test the types
	// For now, just verify the mock server works
	resp, err := http.Get(mock.URL + "/repos/testowner/testrepo/pulls/379")
	if err != nil {
		t.Fatalf("mock server GET: %v", err)
	}
	defer resp.Body.Close()

	var prResp githubPRResponse
	if err := json.NewDecoder(resp.Body).Decode(&prResp); err != nil {
		t.Fatalf("decode mock response: %v", err)
	}
	if prResp.Number != 379 {
		t.Errorf("number = %d, want 379", prResp.Number)
	}
	if prResp.Head.Repo.CloneURL != "https://github.com/forkowner/testrepo.git" {
		t.Errorf("clone URL = %s", prResp.Head.Repo.CloneURL)
	}
}

func TestPrintPR(t *testing.T) {
	info := &PRInfo{
		Number:  379,
		Title:   "Test PR Title",
		Body:    "This is the PR body.",
		State:   "open",
		HTMLURL: "https://github.com/owner/repo/pull/379",
		Head: PRBranch{
			Ref: "feature",
			Repo: PRRepo{
				FullName: "fork/repo",
			},
		},
		Base: PRBranch{
			Ref: "main",
			Repo: PRRepo{
				FullName: "owner/repo",
			},
		},
		User: PRUser{Login: "dev"},
	}

	files := []PRFile{
		{Filename: "main.go", Status: "modified", Additions: 5, Deletions: 3, Changes: 8},
		{Filename: "util.go", Status: "added", Additions: 10, Deletions: 0, Changes: 10},
	}

	diff := "diff --git a/main.go b/main.go\n@@ -1 +1 @@\n-old\n+new\n"

	var buf strings.Builder
	printPR(&buf, info, files, diff)
	output := buf.String()

	if !strings.Contains(output, "PR #379: Test PR Title") {
		t.Errorf("output missing title: %s", output)
	}
	if !strings.Contains(output, "Author: @dev") {
		t.Errorf("output missing author: %s", output)
	}
	if !strings.Contains(output, "fork/repo → owner/repo:main") {
		t.Errorf("output missing branch info: %s", output)
	}
	if !strings.Contains(output, "This is the PR body") {
		t.Errorf("output missing body: %s", output)
	}
	if !strings.Contains(output, "2 files changed, +15 -3") {
		t.Errorf("output missing file summary: %s", output)
	}
	if !strings.Contains(output, "main.go") {
		t.Errorf("output missing file list: %s", output)
	}
	if !strings.Contains(output, "Diff") {
		t.Errorf("output missing diff section: %s", output)
	}
}

func TestPrintPRMinimal(t *testing.T) {
	info := &PRInfo{
		Number:  1,
		Title:   "Minimal",
		State:   "closed",
		HTMLURL: "https://github.com/o/r/pull/1",
		Head: PRBranch{
			Ref:  "fix",
			Repo: PRRepo{FullName: "o/r"},
		},
		Base: PRBranch{
			Ref:  "main",
			Repo: PRRepo{FullName: "o/r"},
		},
		User: PRUser{Login: "u"},
	}

	var buf strings.Builder
	printPR(&buf, info, nil, "")
	output := buf.String()

	if !strings.Contains(output, "PR #1: Minimal") {
		t.Errorf("output missing title: %s", output)
	}
	if strings.Contains(output, "Description") {
		t.Error("output should not have description section when body is empty")
	}
	if strings.Contains(output, "Files Changed") {
		t.Error("output should not have files section when files is nil")
	}
	if strings.Contains(output, "Diff") {
		t.Error("output should not have diff section when diff is empty")
	}
}

func TestHelpText(t *testing.T) {
	if !strings.Contains(help, "github-fetch") {
		t.Error("help text missing github-fetch reference")
	}
	if !strings.Contains(help, "pr") {
		t.Error("help text missing pr command")
	}
	if !strings.Contains(help, "work-on") {
		t.Error("help text missing work-on command")
	}
	if !strings.Contains(help, "push") {
		t.Error("help text missing push command")
	}
	if !strings.Contains(help, "skill install") {
		t.Error("help text missing skill install command")
	}
	if !strings.Contains(help, "skill show") {
		t.Error("help text missing skill show command")
	}
}

func TestHandleNoArgs(t *testing.T) {
	output, err := captureStdout(t, func() error {
		return handle(nil)
	})
	if err != nil {
		t.Fatalf("handle(nil): %v", err)
	}
	if !strings.Contains(output, "Usage:") {
		t.Errorf("no-args output missing Usage: %s", output)
	}
}

func TestHandleHelpShort(t *testing.T) {
	output, err := captureStdout(t, func() error {
		return handle([]string{"-h"})
	})
	if err != nil {
		t.Fatalf("handle(-h): %v", err)
	}
	if !strings.Contains(output, "Usage:") {
		t.Errorf("help output missing Usage: %s", output)
	}
}

func TestHandleHelpLong(t *testing.T) {
	output, err := captureStdout(t, func() error {
		return handle([]string{"--help"})
	})
	if err != nil {
		t.Fatalf("handle(--help): %v", err)
	}
	if !strings.Contains(output, "Usage:") {
		t.Errorf("help output missing Usage: %s", output)
	}
}

func TestHandleSkillShow(t *testing.T) {
	output, err := captureStdout(t, func() error {
		return handle([]string{"skill", "show"})
	})
	if err != nil {
		t.Fatalf("handle(skill show): %v", err)
	}
	if !strings.Contains(output, "github-fetch") {
		t.Errorf("skill show output missing skill name: %s", output)
	}
}

func TestHandleSkillUnknown(t *testing.T) {
	err := handle([]string{"skill", "unknown"})
	if err == nil {
		t.Fatal("expected error for unknown skill sub-command")
	}
	if !strings.Contains(err.Error(), "unknown skill sub-command") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestHandleSkillNoSubcommand(t *testing.T) {
	err := handle([]string{"skill"})
	if err == nil {
		t.Fatal("expected error for skill without sub-command")
	}
	if !strings.Contains(err.Error(), "unknown skill sub-command") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestHandleSkillInstallUnknown(t *testing.T) {
	err := handle([]string{"skill", "unknown-install"})
	if err == nil {
		t.Fatal("expected error for unknown skill sub-command")
	}
	if !strings.Contains(err.Error(), "unknown skill sub-command") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestHandleInstallStandaloneIsUnknown(t *testing.T) {
	err := handle([]string{"install"})
	if err == nil {
		t.Fatal("expected error for standalone install (should be skill install)")
	}
	if !strings.Contains(err.Error(), "unknown command") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestHandleUnknownCommand(t *testing.T) {
	err := handle([]string{"unknown-cmd"})
	if err == nil {
		t.Fatal("expected error for unknown command")
	}
	if !strings.Contains(err.Error(), "unknown command") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestHandlePrWithoutArgs(t *testing.T) {
	err := handle([]string{"pr"})
	if err == nil {
		t.Fatal("expected error for pr without arguments")
	}
	if !strings.Contains(err.Error(), "requires") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestHandleWorkonWithoutArgs(t *testing.T) {
	err := handle([]string{"work-on"})
	if err == nil {
		t.Fatal("expected error for work-on without arguments")
	}
	if !strings.Contains(err.Error(), "work-on requires") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestHandleWorkonOldNameIsUnknown(t *testing.T) {
	err := handle([]string{"workon"})
	if err == nil {
		t.Fatal("expected error for old workon command name")
	}
	if !strings.Contains(err.Error(), "unknown command") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestPushHelp(t *testing.T) {
	output, err := captureStdout(t, func() error {
		return handlePush([]string{"-h"})
	})
	if err != nil {
		t.Fatalf("handlePush(-h): %v", err)
	}
	if !strings.Contains(output, "Usage:") {
		t.Errorf("push help missing Usage: %s", output)
	}
}

func TestHandlePrWithNumberNotInGitRepo(t *testing.T) {
	// Chdir to a temp dir that is not a git repo
	dir := t.TempDir()
	oldDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(oldDir)

	err := handle([]string{"pr", "123"})
	if err == nil {
		t.Fatal("expected error for pr with number when not in git repo")
	}
	if !strings.Contains(err.Error(), "cannot auto-detect") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestHandleWorkonWithInvalidURL(t *testing.T) {
	err := handle([]string{"work-on", "not-a-valid-url"})
	if err == nil {
		t.Fatal("expected error for work-on with invalid URL")
	}
	if !strings.Contains(err.Error(), "invalid PR reference") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestHandlePrWithInvalidURL(t *testing.T) {
	err := handle([]string{"pr", "not-a-valid-url"})
	if err == nil {
		t.Fatal("expected error for pr with invalid URL")
	}
	if !strings.Contains(err.Error(), "invalid PR reference") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestResolvePRFromBranch(t *testing.T) {
	dir := initGitRepo(t, "https://github.com/owner/repo.git")
	oldDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(oldDir)

	// Create a pr-42 branch
	mustRun(t, dir, "git", "checkout", "-b", "pr-42")

	owner, repo, number, err := resolvePRFromBranch()
	if err != nil {
		t.Fatalf("resolvePRFromBranch: %v", err)
	}
	if owner != "owner" {
		t.Errorf("owner = %q, want owner", owner)
	}
	if repo != "repo" {
		t.Errorf("repo = %q, want repo", repo)
	}
	if number != "42" {
		t.Errorf("number = %q, want 42", number)
	}
}

func TestResolvePRFromBranchNotOnPR(t *testing.T) {
	dir := initGitRepo(t, "https://github.com/owner/repo.git")
	oldDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(oldDir)

	_, _, _, err := resolvePRFromBranch()
	if err == nil {
		t.Fatal("expected error when not on PR branch")
	}
	if !strings.Contains(err.Error(), "not on a PR branch") {
		t.Errorf("unexpected error message: %v", err)
	}
}

// Helpers

func initGitRepo(t *testing.T, originURL string) string {
	t.Helper()
	dir := t.TempDir()
	mustRun(t, dir, "git", "init")
	mustRun(t, dir, "git", "config", "user.email", "test@example.com")
	mustRun(t, dir, "git", "config", "user.name", "Test User")
	// Need an initial commit for worktree operations
	writeFile(t, filepath.Join(dir, "README.md"), "# test")
	mustRun(t, dir, "git", "add", "README.md")
	mustRun(t, dir, "git", "commit", "-m", "initial")
	mustRun(t, dir, "git", "remote", "add", "origin", originURL)
	return dir
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("writeFile(%s): %v", path, err)
	}
}

func mustRun(t *testing.T, dir, name string, args ...string) {
	t.Helper()
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
	if err := cmd.Run(); err != nil {
		t.Fatalf("%s %v: %v", name, args, err)
	}
}

func captureStdout(t *testing.T, fn func() error) (string, error) {
	t.Helper()
	oldStdout := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("create stdout pipe: %v", err)
	}

	stdoutCh := make(chan []byte, 1)
	readErrCh := make(chan error, 1)
	go func() {
		data, readErr := io.ReadAll(reader)
		stdoutCh <- data
		readErrCh <- readErr
	}()

	os.Stdout = writer
	runErr := fn()
	os.Stdout = oldStdout
	if err := writer.Close(); err != nil {
		t.Fatalf("close stdout writer: %v", err)
	}
	data := <-stdoutCh
	if err := <-readErrCh; err != nil {
		t.Fatalf("read stdout: %v", err)
	}
	if err := reader.Close(); err != nil {
		t.Fatalf("close stdout reader: %v", err)
	}
	return string(data), runErr
}
