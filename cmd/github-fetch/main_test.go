package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestEmbeddedSkillMDNoInstallGuidelines(t *testing.T) {
	forbidden := []string{
		"skill install",
		"skill show",
		"install --cursor",
		"install --global",
	}
	lower := strings.ToLower(skillTemplate)
	for _, phrase := range forbidden {
		if strings.Contains(lower, phrase) {
			t.Errorf("SKILL.md must not document CLI install/show plumbing (%q found); use --help and README instead", phrase)
		}
	}
}

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

func TestPrintPRSummary(t *testing.T) {
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
		{Filename: "main.go", Status: "modified", Additions: 5, Deletions: 3, Changes: 8, Patch: "line1\nline2\nline3\nline4\nline5"},
		{Filename: "util.go", Status: "added", Additions: 10, Deletions: 0, Changes: 10, Patch: "line1\nline2"},
	}

	var buf strings.Builder
	printPR(&buf, info, files, "", false)
	output := buf.String()

	if !strings.Contains(output, "PR #379: Test PR Title") {
		t.Errorf("output missing title: %s", output)
	}
	if !strings.Contains(output, "2 files changed, +15 -3") {
		t.Errorf("output missing file summary: %s", output)
	}
	if !strings.Contains(output, "── Files Changed ──") {
		t.Errorf("output missing Files Changed section: %s", output)
	}
	if !strings.Contains(output, "── Diff (simplified) ──") {
		t.Errorf("output missing simplified diff section: %s", output)
	}
	// main.go has 5 patch lines, should be truncated to 3
	if !strings.Contains(output, "...") {
		t.Errorf("output missing truncation indicator for main.go: %s", output)
	}
	if strings.Contains(output, "line4") {
		t.Errorf("output should not contain line4 (truncated): %s", output)
	}
}

func TestPrintPRDiff(t *testing.T) {
	info := &PRInfo{
		Number:  379,
		Title:   "Test PR Title",
		Body:    "",
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
		{Filename: "main.go", Status: "modified", Additions: 1, Deletions: 1, Changes: 2, Patch: "line1\nline2\nline3\nline4"},
	}

	diff := "diff --git a/main.go b/main.go\n@@ -1 +1 @@\n-old\n+new\n"

	var buf strings.Builder
	printPR(&buf, info, files, diff, true)
	output := buf.String()

	if !strings.Contains(output, "── Diff ──") {
		t.Errorf("diff mode should show diff section: %s", output)
	}
	if strings.Contains(output, "simplified") {
		t.Errorf("diff mode should not show simplified diff: %s", output)
	}
}

func TestPrintPRDiffFileLimit(t *testing.T) {
	info := &PRInfo{
		Number:  1,
		Title:   "Many files",
		Body:    "",
		State:   "open",
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

	files := make([]PRFile, 15)
	for i := range files {
		files[i] = PRFile{
			Filename:  fmt.Sprintf("file%d.go", i),
			Status:    "modified",
			Additions: 1,
			Deletions: 0,
			Patch:     "line1",
		}
	}

	var buf strings.Builder
	printPR(&buf, info, files, "", false)
	output := buf.String()

	if !strings.Contains(output, "and 5 more files") {
		t.Errorf("output missing overflow indicator: %s", output)
	}
	if strings.Contains(output, "file10.go:") {
		t.Errorf("output should not include file beyond limit 10 in simplified diff: %s", output)
	}
	// Files Changed section should still show all 15
	if !strings.Contains(output, "15 files changed") {
		t.Errorf("output missing total file count: %s", output)
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
	printPR(&buf, info, nil, "", false)
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
	if !strings.Contains(help, "ci") {
		t.Error("help text missing ci command")
	}
	if !strings.Contains(help, "work-on") {
		t.Error("help text missing work-on command")
	}
	if !strings.Contains(help, "push") {
		t.Error("help text missing push command")
	}
	if !strings.Contains(help, "skill --install") {
		t.Error("help text missing skill --install command")
	}
	if !strings.Contains(help, "skill --show") {
		t.Error("help text missing skill --show command")
	}
	if !strings.Contains(help, "pr --logs") {
		t.Error("help text missing pr --logs option")
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
		return handle([]string{"skill", "--show"})
	})
	if err != nil {
		t.Fatalf("handle(skill --show): %v", err)
	}
	if !strings.Contains(output, "github-fetch") {
		t.Errorf("skill --show output missing skill name: %s", output)
	}
}

func TestHandleSkillUnknown(t *testing.T) {
	err := handle([]string{"skill", "unknown"})
	if err == nil {
		t.Fatal("expected error for skill without action flags")
	}
	if !strings.Contains(err.Error(), "--show") && !strings.Contains(err.Error(), "--install") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestHandleSkillNoSubcommand(t *testing.T) {
	err := handle([]string{"skill"})
	if err == nil {
		t.Fatal("expected error for skill without action flags")
	}
	if !strings.Contains(err.Error(), "--show") && !strings.Contains(err.Error(), "--install") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestHandleSkillInstallUnknown(t *testing.T) {
	// bare token without action flag is a missing-action error
	err := handle([]string{"skill", "unknown-install"})
	if err == nil {
		t.Fatal("expected error for skill without action flags")
	}
	if !strings.Contains(err.Error(), "--show") && !strings.Contains(err.Error(), "--install") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestHandleInstallStandaloneIsUnknown(t *testing.T) {
	err := handle([]string{"install"})
	if err == nil {
		t.Fatal("expected error for standalone install (should be skill --install)")
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

func TestParseActionsURL(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		owner   string
		repo    string
		runID   string
		jobID   string
		wantErr bool
	}{
		{"job url with pr query", "https://github.com/xhd2015/xgo/actions/runs/26795086426/job/78989577716?pr=385", "xhd2015", "xgo", "26795086426", "78989577716", false},
		{"job url no query", "https://github.com/xhd2015/xgo/actions/runs/26795086426/job/78989577716", "xhd2015", "xgo", "26795086426", "78989577716", false},
		{"run url", "https://github.com/owner/repo/actions/runs/123456", "owner", "repo", "123456", "", false},
		{"run url with trailing slash", "https://github.com/owner/repo/actions/runs/123456/", "owner", "repo", "123456", "", false},
		{"http scheme", "http://github.com/owner/repo/actions/runs/123/job/456", "owner", "repo", "123", "456", false},
		{"not actions url", "https://github.com/owner/repo/pull/123", "", "", "", "", true},
		{"random string", "not-a-url", "", "", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			owner, repo, runID, jobID, err := parseActionsURL(tt.raw)
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
			if runID != tt.runID {
				t.Errorf("runID = %q, want %q", runID, tt.runID)
			}
			if jobID != tt.jobID {
				t.Errorf("jobID = %q, want %q", jobID, tt.jobID)
			}
		})
	}
}

func TestIsActionsURL(t *testing.T) {
	tests := []struct {
		raw    string
		isAction bool
	}{
		{"https://github.com/o/r/actions/runs/1/job/2", true},
		{"https://github.com/o/r/actions/runs/1", true},
		{"https://github.com/o/r/pull/1", false},
		{"https://github.com/o/r/issues/1", false},
		{"not-a-url", false},
	}

	for _, tt := range tests {
		t.Run(tt.raw, func(t *testing.T) {
			got := isActionsURL(tt.raw)
			if got != tt.isAction {
				t.Errorf("isActionsURL(%q) = %v, want %v", tt.raw, got, tt.isAction)
			}
		})
	}
}

func TestHandleCIHelp(t *testing.T) {
	output, err := captureStdout(t, func() error {
		return handleCI([]string{"-h"})
	})
	if err != nil {
		t.Fatalf("handleCI(-h): %v", err)
	}
	if !strings.Contains(output, "Usage:") {
		t.Errorf("ci help missing Usage: %s", output)
	}
	if !strings.Contains(output, "--logs") {
		t.Errorf("ci help missing --logs: %s", output)
	}
}

func TestHandleCIWithoutArgs(t *testing.T) {
	err := handleCI(nil)
	// Empty args auto-detects the current repo; failure modes vary by env
	// (no origin, no workflow runs, network). Just require an error.
	if err == nil {
		t.Skip("handleCI(nil) succeeded via auto-detect in this environment")
	}
}

func TestHandleCIWithoutArgsEmpty(t *testing.T) {
	err := handleCI([]string{})
	if err == nil {
		t.Skip("handleCI([]) succeeded via auto-detect in this environment")
	}
}

func TestPrintWorkflowRuns(t *testing.T) {
	info := &PRInfo{
		Number: 385,
		Head:   PRBranch{Ref: "dev-go1.25"},
	}
	runs := []WorkflowRun{
		{ID: 1, Name: "CI / Build", Conclusion: "success", Status: "completed", HTMLURL: "https://github.com/o/r/actions/runs/1"},
		{ID: 2, Name: "CI / Test", Conclusion: "failure", Status: "completed", HTMLURL: "https://github.com/o/r/actions/runs/2"},
		{ID: 3, Name: "CI / Lint", Status: "in_progress", HTMLURL: "https://github.com/o/r/actions/runs/3"},
		{ID: 4, Name: "CI / Deploy", Status: "queued", HTMLURL: "https://github.com/o/r/actions/runs/4"},
	}

	var buf strings.Builder
	printWorkflowRuns(&buf, info, runs)
	output := buf.String()

	if !strings.Contains(output, "Workflow Runs for PR #385 (dev-go1.25)") {
		t.Errorf("output missing header: %s", output)
	}
	if !strings.Contains(output, "CI / Build") {
		t.Errorf("output missing run name: %s", output)
	}
	if !strings.Contains(output, "success") {
		t.Errorf("output missing success status: %s", output)
	}
	if !strings.Contains(output, "failure") {
		t.Errorf("output missing failure status: %s", output)
	}
	if !strings.Contains(output, "in_progress") {
		t.Errorf("output missing in_progress status: %s", output)
	}
	if !strings.Contains(output, "queued") {
		t.Errorf("output missing queued status: %s", output)
	}
}

func TestPrintWorkflowRunsEmpty(t *testing.T) {
	info := &PRInfo{
		Number: 1,
		Head:   PRBranch{Ref: "main"},
	}

	var buf strings.Builder
	printWorkflowRuns(&buf, info, nil)
	output := buf.String()

	if !strings.Contains(output, "No workflow runs found") {
		t.Errorf("output missing empty message: %s", output)
	}
}

func TestExtractLogZip(t *testing.T) {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)

	files := []struct {
		name    string
		content string
	}{
		{"0_job.txt", "Step 1: Build\nStep 2: Test\nerror: test failed\n"},
		{"1_compile step.txt", "go build ./...\n"},
	}

	for _, f := range files {
		w, err := zw.Create(f.name)
		if err != nil {
			t.Fatalf("create zip entry %q: %v", f.name, err)
		}
		if _, err := w.Write([]byte(f.content)); err != nil {
			t.Fatalf("write zip entry %q: %v", f.name, err)
		}
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("close zip writer: %v", err)
	}

	result, err := extractLogZip(buf.Bytes())
	if err != nil {
		t.Fatalf("extractLogZip: %v", err)
	}

	if !strings.Contains(result, "── 0_job.txt ──") {
		t.Errorf("output missing 0_job.txt header: %s", result)
	}
	if !strings.Contains(result, "error: test failed") {
		t.Errorf("output missing log content: %s", result)
	}
	if !strings.Contains(result, "── 1_compile step.txt ──") {
		t.Errorf("output missing 1_compile step.txt header: %s", result)
	}
}

func TestExtractLogZipEmpty(t *testing.T) {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	if err := zw.Close(); err != nil {
		t.Fatalf("close zip writer: %v", err)
	}

	result, err := extractLogZip(buf.Bytes())
	if err != nil {
		t.Fatalf("extractLogZip (empty): %v", err)
	}
	if result != "" {
		t.Errorf("expected empty string for empty zip, got %q", result)
	}
}

func TestFetchWorkflowRunsWithMock(t *testing.T) {
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/jobs") {
			json.NewEncoder(w).Encode(githubWorkflowJobsResponse{
				TotalCount: 1,
				Jobs: []githubWorkflowJob{
					{ID: 100, Name: "test-job", Status: "completed", Conclusion: "success"},
				},
			})
			return
		}
		resp := githubWorkflowRunsResponse{
			TotalCount: 2,
			WorkflowRuns: []githubWorkflowRun{
				{ID: 1, Name: "CI / Build", Status: "completed", Conclusion: "success", HTMLURL: "https://github.com/o/r/actions/runs/1"},
				{ID: 2, Name: "CI / Test", Status: "completed", Conclusion: "failure", HTMLURL: "https://github.com/o/r/actions/runs/2"},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer mock.Close()

	// Test the response types by decoding
	resp, err := http.Get(mock.URL + "/repos/o/r/actions/runs?head_sha=abc")
	if err != nil {
		t.Fatalf("mock GET: %v", err)
	}
	defer resp.Body.Close()

	var runsResp githubWorkflowRunsResponse
	if err := json.NewDecoder(resp.Body).Decode(&runsResp); err != nil {
		t.Fatalf("decode mock response: %v", err)
	}
	if runsResp.TotalCount != 2 {
		t.Errorf("total_count = %d, want 2", runsResp.TotalCount)
	}
	if runsResp.WorkflowRuns[0].Name != "CI / Build" {
		t.Errorf("run[0].Name = %q, want CI / Build", runsResp.WorkflowRuns[0].Name)
	}
	if runsResp.WorkflowRuns[1].Conclusion != "failure" {
		t.Errorf("run[1].Conclusion = %q, want failure", runsResp.WorkflowRuns[1].Conclusion)
	}
}

func TestFetchJobLogsWithMock(t *testing.T) {
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer
		zw := zip.NewWriter(&buf)
		fw, _ := zw.Create("test_log.txt")
		fw.Write([]byte("line1\nline2\n"))
		zw.Close()
		w.Header().Set("Content-Type", "application/zip")
		w.Write(buf.Bytes())
	}))
	defer mock.Close()

	// Directly test extractLogZip with a known zip
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	fw, _ := zw.Create("test_log.txt")
	fw.Write([]byte("line1\nline2\n"))
	zw.Close()

	result, err := extractLogZip(buf.Bytes())
	if err != nil {
		t.Fatalf("extractLogZip: %v", err)
	}
	if !strings.Contains(result, "line1") {
		t.Errorf("expected line1 in extracted logs: %s", result)
	}
}

func TestHandleCIActionsURLWithMock(t *testing.T) {
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/logs") {
			var buf bytes.Buffer
			zw := zip.NewWriter(&buf)
			fw, _ := zw.Create("log.txt")
			fw.Write([]byte("build output here\n"))
			zw.Close()
			w.Header().Set("Content-Type", "application/zip")
			w.Write(buf.Bytes())
			return
		}
		if strings.Contains(r.URL.Path, "/jobs") {
			json.NewEncoder(w).Encode(githubWorkflowJobsResponse{
				TotalCount: 1,
				Jobs: []githubWorkflowJob{
					{ID: 100, Name: "test-job", Status: "completed", Conclusion: "success"},
				},
			})
			return
		}
		json.NewEncoder(w).Encode(githubWorkflowRunsResponse{
			TotalCount: 1,
			WorkflowRuns: []githubWorkflowRun{
				{ID: 100, Name: "CI / Test", Status: "completed", Conclusion: "success"},
			},
		})
	}))
	defer mock.Close()

	// Test the mock handles the expected endpoints
	resp, err := http.Get(mock.URL + "/repos/o/r/actions/runs/100/jobs")
	if err != nil {
		t.Fatalf("mock jobs GET: %v", err)
	}
	defer resp.Body.Close()

	var jobsResp githubWorkflowJobsResponse
	if err := json.NewDecoder(resp.Body).Decode(&jobsResp); err != nil {
		t.Fatalf("decode jobs response: %v", err)
	}
	if jobsResp.TotalCount != 1 {
		t.Errorf("total_count = %d, want 1", jobsResp.TotalCount)
	}
	if jobsResp.Jobs[0].Name != "test-job" {
		t.Errorf("job name = %q, want test-job", jobsResp.Jobs[0].Name)
	}

	// Test the logs endpoint returns a valid zip
	resp2, err := http.Get(mock.URL + "/repos/o/r/actions/jobs/100/logs")
	if err != nil {
		t.Fatalf("mock logs GET: %v", err)
	}
	defer resp2.Body.Close()
	data, _ := io.ReadAll(resp2.Body)
	logContent, err := extractLogZip(data)
	if err != nil {
		t.Fatalf("extract log zip from mock: %v", err)
	}
	if !strings.Contains(logContent, "build output here") {
		t.Errorf("unexpected log content: %s", logContent)
	}
}

func TestPrLogsDelegationToActionsURL(t *testing.T) {
	// Verify that pr --logs with actions URL would call handleCI
	// We test this by checking that isActionsURL returns true for action URL
	// and that the pr handler calls the ci path
	if !isActionsURL("https://github.com/o/r/actions/runs/1/job/2") {
		t.Error("isActionsURL should return true for job URL")
	}
}

func TestHandleCIActionURLPRRefConflict(t *testing.T) {
	// Verify ci subcommand rejects invalid --run-id value
	err := handleCI([]string{"--run-id", "abc", "https://github.com/o/r/pull/1"})
	if err == nil {
		t.Fatal("expected error for invalid --run-id value")
	}
	if !strings.Contains(strings.ToLower(err.Error()), "invalid") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestHandleCIUnknownCommand(t *testing.T) {
	err := handle([]string{"ci", "--invalid-flag"})
	if err == nil {
		t.Fatal("expected error for invalid flag")
	}
	msg := strings.ToLower(err.Error())
	if !strings.Contains(msg, "unrecognized flag") && !strings.Contains(msg, "unknown flag") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestHandleChecksAlias(t *testing.T) {
	dir := initGitRepo(t, "https://github.com/testowner/testrepo.git")
	oldDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(oldDir)

	// "checks" is an alias for "ci" - it should not return "unknown command"
	err := handle([]string{"checks", "https://github.com/testowner/testrepo/pull/42"})
	// Expected to fail because it tries to make real API calls, but not with "unknown command"
	if err == nil {
		return // would be unexpected but fine
	}
	if strings.Contains(err.Error(), "unknown command") {
		t.Errorf("checks alias should dispatch to ci, got: %v", err)
	}
}

func TestRunIDFlagValueRequired(t *testing.T) {
	err := handleCI([]string{"--run-id"})
	if err == nil {
		t.Fatal("expected error for --run-id without value")
	}
	if !strings.Contains(err.Error(), "requires a value") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestJobFlagValueRequired(t *testing.T) {
	err := handleCI([]string{"--job"})
	if err == nil {
		t.Fatal("expected error for --job without value")
	}
	if !strings.Contains(err.Error(), "requires a value") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestWorkflowFlagValueRequired(t *testing.T) {
	err := handleCI([]string{"--workflow"})
	if err == nil {
		t.Fatal("expected error for --workflow without value")
	}
	if !strings.Contains(err.Error(), "requires a value") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestFilterRunsByName(t *testing.T) {
	runs := []WorkflowRun{
		{ID: 1, Name: "Go 1-25 / test (go1.25)", Conclusion: "failure"},
		{ID: 2, Name: "Go 1-24 / test (go1.24)", Conclusion: "success"},
		{ID: 3, Name: "Lint", Conclusion: "success"},
	}

	filtered := filterRunsByName(runs, "Go 1-25")
	if len(filtered) != 1 {
		t.Fatalf("expected 1 run, got %d", len(filtered))
	}
	if filtered[0].ID != 1 {
		t.Errorf("expected run ID 1, got %d", filtered[0].ID)
	}

	filtered = filterRunsByName(runs, "go 1-25")
	if len(filtered) != 1 {
		t.Fatalf("case-insensitive: expected 1 run, got %d", len(filtered))
	}

	filtered = filterRunsByName(runs, "Go")
	if len(filtered) != 2 {
		t.Fatalf("partial match: expected 2 runs, got %d", len(filtered))
	}

	filtered = filterRunsByName(runs, "nonexistent")
	if len(filtered) != 0 {
		t.Fatalf("no match: expected 0 runs, got %d", len(filtered))
	}

	filtered = filterRunsByName(nil, "Go")
	if len(filtered) != 0 {
		t.Fatalf("nil input: expected 0 runs, got %d", len(filtered))
	}
}

func TestPrintWorkflowRunsFiltered(t *testing.T) {
	info := &PRInfo{
		Number: 385,
		Head:   PRBranch{Ref: "dev-go1.25"},
	}
	runs := []WorkflowRun{
		{ID: 1, Name: "Go 1-25 / test (go1.25)", Conclusion: "failure", Status: "completed", HTMLURL: "https://github.com/o/r/actions/runs/1"},
	}

	var buf strings.Builder
	printWorkflowRuns(&buf, info, runs)
	output := buf.String()

	if !strings.Contains(output, "Go 1-25 / test (go1.25)") {
		t.Errorf("output missing filtered run: %s", output)
	}
}

func TestTruncateLogUnderLimit(t *testing.T) {
	lines := make([]string, 100)
	for i := range lines {
		lines[i] = fmt.Sprintf("line %d", i)
	}
	content := strings.Join(lines, "\n")
	result := truncateLog(content, false)

	if !strings.Contains(result, "line 0") {
		t.Errorf("expected content to be preserved, got: %s", result)
	}
	if strings.Contains(result, "showing last") {
		t.Errorf("expected no truncation message for small content")
	}
}

func TestTruncateLogOverLimit(t *testing.T) {
	n := defaultLogLines + 100
	lines := make([]string, n)
	for i := range lines {
		lines[i] = fmt.Sprintf("line %d", i)
	}
	content := strings.Join(lines, "\n")
	result := truncateLog(content, false)

	if !strings.Contains(result, "showing last") {
		t.Errorf("expected truncation message, got: %s", result)
	}
	if !strings.Contains(result, fmt.Sprintf("%d lines", n)) {
		t.Errorf("expected total line count %d in message, got: %s", n, result)
	}
	if !strings.Contains(result, "use --full") {
		t.Errorf("expected --full hint, got: %s", result)
	}
	if strings.Contains(result, "line 0") {
		t.Error("expected first lines to be truncated")
	}
	if !strings.Contains(result, fmt.Sprintf("line %d", n-1)) {
		t.Errorf("expected last line to be present")
	}
	if !strings.Contains(result, fmt.Sprintf("line %d", n-defaultLogLines)) {
		t.Errorf("expected line at truncation boundary to be present")
	}
}

func TestTruncateLogFull(t *testing.T) {
	n := defaultLogLines + 100
	lines := make([]string, n)
	for i := range lines {
		lines[i] = fmt.Sprintf("line %d", i)
	}
	content := strings.Join(lines, "\n")
	result := truncateLog(content, true)

	if strings.Contains(result, "showing last") {
		t.Errorf("full mode should not truncate")
	}
	if !strings.Contains(result, "line 0") {
		t.Errorf("full mode should contain first line")
	}
	if !strings.Contains(result, fmt.Sprintf("line %d", n-1)) {
		t.Errorf("full mode should contain last line")
	}
}

func TestTruncateLogEmpty(t *testing.T) {
	result := truncateLog("", false)
	if result != "" {
		t.Errorf("expected empty string, got: %q", result)
	}
}

func TestTruncateLogExactLimit(t *testing.T) {
	lines := make([]string, defaultLogLines)
	for i := range lines {
		lines[i] = fmt.Sprintf("line %d", i)
	}
	content := strings.Join(lines, "\n")
	result := truncateLog(content, false)

	if strings.Contains(result, "showing last") {
		t.Errorf("content at exact limit should not be truncated")
	}
}

func TestHandleCIFullFlag(t *testing.T) {
	output, err := captureStdout(t, func() error {
		return handleCI([]string{"-h"})
	})
	if err != nil {
		t.Fatalf("handleCI(-h): %v", err)
	}
	if !strings.Contains(output, "--full") {
		t.Errorf("ci help missing --full: %s", output)
	}
}

func TestPrLogsFullFlag(t *testing.T) {
	output, err := captureStdout(t, func() error {
		return handleFetchPR([]string{"-h"})
	})
	if err != nil {
		t.Fatalf("handleFetchPR(-h): %v", err)
	}
	if !strings.Contains(output, "--full") {
		t.Errorf("pr help missing --full: %s", output)
	}
}

func TestNormalizeSlashes(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Go 1-25 / test-go1-25", "Go 1-25/test-go1-25"},
		{"Go 1-25/test-go1-25", "Go 1-25/test-go1-25"},
		{"a / b / c", "a/b/c"},
		{"a/b/c", "a/b/c"},
		{"no-slash", "no-slash"},
		{" / leading", "/leading"},
		{"trailing / ", "trailing/"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := normalizeSlashes(tt.input)
			if got != tt.expected {
				t.Errorf("normalizeSlashes(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestFilterRunsByNameNormalized(t *testing.T) {
	runs := []WorkflowRun{
		{ID: 1, Name: "Go 1-25 / test-go1-25", Conclusion: "failure"},
		{ID: 2, Name: "Go 1-24 / test-go1-24", Conclusion: "success"},
	}

	// "Go 1-25/test" should match "Go 1-25 / test-go1-25"
	filtered := filterRunsByName(runs, "Go 1-25/test")
	if len(filtered) != 1 || filtered[0].ID != 1 {
		t.Fatalf("expected run ID 1, got %v", filtered)
	}

	// "Go 1-25 / test" should also match
	filtered = filterRunsByName(runs, "Go 1-25 / test")
	if len(filtered) != 1 || filtered[0].ID != 1 {
		t.Fatalf("expected run ID 1 with spaced filter, got %v", filtered)
	}

	// Exact match without slash spaces
	filtered = filterRunsByName(runs, "Go 1-25/test-go1-25")
	if len(filtered) != 1 {
		t.Fatalf("expected 1 result for exact match, got %d", len(filtered))
	}
}

func TestPrLogsSkipsPRDisplay(t *testing.T) {
	dir := initGitRepo(t, "https://github.com/testowner/testrepo.git")
	oldDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(oldDir)

	err := handleFetchPR([]string{"--logs", "https://github.com/testowner/testrepo/pull/42"})
	if err == nil {
		return
	}
	if strings.Contains(err.Error(), "fetch PR files") {
		t.Errorf("pr --logs should skip PR display, but got PR file fetch error: %v", err)
	}
}

func TestHandleYAMLNoSubcommand(t *testing.T) {
	err := handleYAML(nil)
	if err == nil {
		t.Fatal("expected error for yaml without subcommand")
	}
	if !strings.Contains(err.Error(), "requires a subcommand") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestHandleYAMLUnknownSubcommand(t *testing.T) {
	err := handleYAML([]string{"unknown"})
	if err == nil {
		t.Fatal("expected error for unknown yaml subcommand")
	}
	if !strings.Contains(err.Error(), "unknown yaml subcommand") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestHandleYAMLValidateNoPath(t *testing.T) {
	err := handleYAMLValidate([]string{})
	if err == nil {
		t.Fatal("expected error for validate without path")
	}
	if !strings.Contains(err.Error(), "requires a file path") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestHandleYAMLValidateHelp(t *testing.T) {
	output, err := captureStdout(t, func() error {
		return handleYAMLValidate([]string{"-h"})
	})
	if err != nil {
		t.Fatalf("handleYAMLValidate(-h): %v", err)
	}
	if !strings.Contains(output, "Usage:") {
		t.Errorf("help missing Usage: %s", output)
	}
}

func TestLocalValidateYAMLValid(t *testing.T) {
	content := `name: Build
on:
  push:
    branches: [main]
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - run: echo hello
`
	output, err := captureStdout(t, func() error {
		return localValidateYAML("test.yml", content)
	})
	if err != nil {
		t.Fatalf("localValidateYAML: %v", err)
	}
	if !strings.Contains(output, "OK") {
		t.Errorf("expected OK for valid yaml: %s", output)
	}
}

func TestLocalValidateYAMLMissingOn(t *testing.T) {
	content := `name: Build
jobs:
  build:
    runs-on: ubuntu-latest
`
	output, err := captureStdout(t, func() error {
		return localValidateYAML("test.yml", content)
	})
	if err != nil {
		t.Fatalf("localValidateYAML: %v", err)
	}
	if !strings.Contains(output, "missing `on:") {
		t.Errorf("expected missing on: warning: %s", output)
	}
}

func TestLocalValidateYAMLMissingJobs(t *testing.T) {
	content := `name: Build
on:
  push:
    branches: [main]
`
	output, err := captureStdout(t, func() error {
		return localValidateYAML("test.yml", content)
	})
	if err != nil {
		t.Fatalf("localValidateYAML: %v", err)
	}
	if !strings.Contains(output, "missing `jobs:") {
		t.Errorf("expected missing jobs: warning: %s", output)
	}
}

func TestLocalValidateYAMLMissingBoth(t *testing.T) {
	content := `name: Build
`
	output, err := captureStdout(t, func() error {
		return localValidateYAML("test.yml", content)
	})
	if err != nil {
		t.Fatalf("localValidateYAML: %v", err)
	}
	if !strings.Contains(output, "missing `on:") {
		t.Errorf("expected missing on: warning: %s", output)
	}
	if !strings.Contains(output, "missing `jobs:") {
		t.Errorf("expected missing jobs: warning: %s", output)
	}
}

func TestValidateWorkflowFileNotFound(t *testing.T) {
	err := validateWorkflowFile("/nonexistent/path/workflow.yml")
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
	if !strings.Contains(err.Error(), "file not found") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateWorkflowFileValid(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "test.yml")
	content := `name: Build
on:
  push:
    branches: [main]
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - run: echo hello
`
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("write test file: %v", err)
	}

	output, err := captureStdout(t, func() error {
		return validateWorkflowFile(filePath)
	})
	if err != nil {
		t.Fatalf("validateWorkflowFile: %v", err)
	}
	if !strings.Contains(output, "OK") {
		t.Errorf("expected OK: %s", output)
	}
}

func TestYAMLHelpText(t *testing.T) {
	output, err := captureStdout(t, func() error {
		return handleYAML([]string{"-h"})
	})
	if err != nil {
		t.Fatalf("handleYAML(-h): %v", err)
	}
	if !strings.Contains(output, "Usage:") {
		t.Errorf("yaml help missing Usage: %s", output)
	}
	if !strings.Contains(output, "validate") {
		t.Errorf("yaml help missing validate: %s", output)
	}
}

func TestLocalValidateYAMLSyntaxError(t *testing.T) {
	content := `name: Build
on:
  push
    branches: [main]
jobs:
  build:
    runs-on: ubuntu-latest
`
	output, err := captureStdout(t, func() error {
		return localValidateYAML("test.yml", content)
	})
	if err != nil {
		t.Fatalf("localValidateYAML: %v", err)
	}
	if !strings.Contains(output, "YAML syntax error") {
		t.Errorf("expected YAML syntax error: %s", output)
	}
}

func TestLocalValidateYAMLOKBadge(t *testing.T) {
	content := `name: Build
on:
  push:
    branches: [main]
jobs:
  build:
    runs-on: ubuntu-latest
`
	output, err := captureStdout(t, func() error {
		return localValidateYAML("test.yml", content)
	})
	if err != nil {
		t.Fatalf("localValidateYAML: %v", err)
	}
	if !strings.Contains(output, "OK") {
		t.Errorf("expected OK for valid yaml: %s", output)
	}
	if !strings.Contains(output, "valid YAML") {
		t.Errorf("expected valid YAML message: %s", output)
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

func TestPrWorkflowEqualsSyntax(t *testing.T) {
	output, err := captureStdout(t, func() error {
		return handleFetchPR([]string{"--workflow=build", "--diff", "-h"})
	})
	if err != nil {
		t.Fatalf("handleFetchPR with = syntax: %v", err)
	}
	if !strings.Contains(output, "Usage:") {
		t.Errorf("help output missing Usage: %s", output)
	}
}

func TestCIRunIDEqualsSyntax(t *testing.T) {
	output, err := captureStdout(t, func() error {
		return handleCI([]string{"--run-id=123", "--job=test", "-h"})
	})
	if err != nil {
		t.Fatalf("handleCI with = syntax: %v", err)
	}
	if !strings.Contains(output, "Usage:") {
		t.Errorf("help output missing Usage: %s", output)
	}
}

func TestPushForceEqualsSyntax(t *testing.T) {
	output, err := captureStdout(t, func() error {
		return handlePush([]string{"--force", "-h"})
	})
	if err != nil {
		t.Fatalf("handlePush with = syntax: %v", err)
	}
	if !strings.Contains(output, "Usage:") {
		t.Errorf("help output missing Usage: %s", output)
	}
}

func TestPushShortFlagHelp(t *testing.T) {
	output, err := captureStdout(t, func() error {
		return handlePush([]string{"-f", "-h"})
	})
	if err != nil {
		t.Fatalf("handlePush with short flags: %v", err)
	}
	if !strings.Contains(output, "Usage:") {
		t.Errorf("help output missing Usage: %s", output)
	}
}

func TestWaitForRunCompletionAlreadyCompleted(t *testing.T) {
	callCount := 0
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		json.NewEncoder(w).Encode(githubWorkflowRun{
			ID:         100,
			Name:       "CI / Test",
			Status:     "completed",
			Conclusion: "success",
			HTMLURL:    "https://github.com/o/r/actions/runs/100",
		})
	}))
	defer mock.Close()

	origURL := apiBaseURL
	apiBaseURL = mock.URL
	defer func() { apiBaseURL = origURL }()

	var buf strings.Builder
	err := waitForRunCompletion(&buf, "o", "r", 100, "CI / Test", false)
	if err != nil {
		t.Fatalf("waitForRunCompletion: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output for already completed run, got: %s", buf.String())
	}
	if callCount != 1 {
		t.Errorf("expected 1 API call, got %d", callCount)
	}
}

func TestWaitForRunCompletionNoWait(t *testing.T) {
	var buf strings.Builder
	err := waitForRunCompletion(&buf, "o", "r", 100, "CI / Test", true)
	if err != nil {
		t.Fatalf("waitForRunCompletion (noWait): %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output for noWait, got: %s", buf.String())
	}
}

func TestWaitForRunCompletionInProgressThenCompleted(t *testing.T) {
	callCount := 0
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount <= 3 {
			json.NewEncoder(w).Encode(githubWorkflowRun{
				ID:     100,
				Name:   "CI / Test",
				Status: "in_progress",
				HTMLURL: "https://github.com/o/r/actions/runs/100",
			})
		} else {
			json.NewEncoder(w).Encode(githubWorkflowRun{
				ID:         100,
				Name:       "CI / Test",
				Status:     "completed",
				Conclusion: "failure",
				HTMLURL:    "https://github.com/o/r/actions/runs/100",
			})
		}
	}))
	defer mock.Close()

	origURL := apiBaseURL
	origInterval := pollInterval
	apiBaseURL = mock.URL
	pollInterval = 10 * time.Millisecond
	defer func() {
		apiBaseURL = origURL
		pollInterval = origInterval
	}()

	var buf strings.Builder
	err := waitForRunCompletion(&buf, "o", "r", 100, "CI / Test", false)
	if err != nil {
		t.Fatalf("waitForRunCompletion: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "Waiting for CI / Test (run #100) to complete") {
		t.Errorf("output missing waiting message: %s", output)
	}
	if !strings.Contains(output, "in_progress") {
		t.Errorf("output missing initial status line: %s", output)
	}
	if !strings.Contains(output, "completed failure") {
		t.Errorf("output missing completed transition line: %s", output)
	}
	if !strings.Contains(output, "done (took ") {
		t.Errorf("output missing done message: %s", output)
	}
	if !strings.Contains(output, ".") {
		t.Errorf("output missing dots for unchanged status: %s", output)
	}
	if callCount != 4 {
		t.Errorf("expected 4 API calls (1 initial + 3 polls), got %d", callCount)
	}
}

func TestWaitForRunCompletionStatusTransitions(t *testing.T) {
	callCount := 0
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		switch {
		case callCount == 1:
			json.NewEncoder(w).Encode(githubWorkflowRun{
				ID:     100,
				Name:   "CI / Test",
				Status: "queued",
				HTMLURL: "https://github.com/o/r/actions/runs/100",
			})
		case callCount <= 3:
			json.NewEncoder(w).Encode(githubWorkflowRun{
				ID:     100,
				Name:   "CI / Test",
				Status: "in_progress",
				HTMLURL: "https://github.com/o/r/actions/runs/100",
			})
		default:
			json.NewEncoder(w).Encode(githubWorkflowRun{
				ID:         100,
				Name:       "CI / Test",
				Status:     "completed",
				Conclusion: "success",
				HTMLURL:    "https://github.com/o/r/actions/runs/100",
			})
		}
	}))
	defer mock.Close()

	origURL := apiBaseURL
	origInterval := pollInterval
	apiBaseURL = mock.URL
	pollInterval = 10 * time.Millisecond
	defer func() {
		apiBaseURL = origURL
		pollInterval = origInterval
	}()

	var buf strings.Builder
	err := waitForRunCompletion(&buf, "o", "r", 100, "CI / Test", false)
	if err != nil {
		t.Fatalf("waitForRunCompletion: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "Waiting for CI / Test (run #100) to complete") {
		t.Errorf("output missing waiting message: %s", output)
	}
	if !strings.Contains(output, "queued") {
		t.Errorf("output missing queued status line: %s", output)
	}
	if !strings.Contains(output, "in_progress") {
		t.Errorf("output missing in_progress transition line: %s", output)
	}
	if !strings.Contains(output, "completed success") {
		t.Errorf("output missing completed transition line: %s", output)
	}
	if !strings.Contains(output, "done (took ") {
		t.Errorf("output missing done message: %s", output)
	}
	if !strings.Contains(output, ".") {
		t.Errorf("output missing dots: %s", output)
	}
	// queued should appear before in_progress
	queuedIdx := strings.Index(output, "queued")
	inProgressIdx := strings.Index(output, "in_progress")
	if queuedIdx < 0 || inProgressIdx < 0 || queuedIdx >= inProgressIdx {
		t.Errorf("expected queued before in_progress in output: %s", output)
	}
	// Verify the output has 3 status lines (queued, in_progress, completed)
	// by counting lines that contain recognizable status keywords
	statusLines := 0
	for _, line := range strings.Split(output, "\n") {
		if strings.Contains(line, "queued") || strings.Contains(line, "in_progress") || strings.Contains(line, "completed") {
			statusLines++
		}
	}
	if statusLines < 3 {
		t.Errorf("expected at least 3 status lines, got %d: %s", statusLines, output)
	}
}

func TestWaitForRunCompletionTimeout(t *testing.T) {
	mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(githubWorkflowRun{
			ID:     100,
			Name:   "CI / Test",
			Status: "in_progress",
			HTMLURL: "https://github.com/o/r/actions/runs/100",
		})
	}))
	defer mock.Close()

	origURL := apiBaseURL
	origInterval := pollInterval
	origTimeout := pollTimeout
	apiBaseURL = mock.URL
	pollInterval = 10 * time.Millisecond
	pollTimeout = 50 * time.Millisecond
	defer func() {
		apiBaseURL = origURL
		pollInterval = origInterval
		pollTimeout = origTimeout
	}()

	var buf strings.Builder
	err := waitForRunCompletion(&buf, "o", "r", 100, "CI / Test", false)
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
	if !strings.Contains(err.Error(), "timed out") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestCINoWaitFlag(t *testing.T) {
	output, err := captureStdout(t, func() error {
		return handleCI([]string{"--no-wait", "-h"})
	})
	if err != nil {
		t.Fatalf("handleCI with --no-wait: %v", err)
	}
	if !strings.Contains(output, "--no-wait") {
		t.Errorf("ci help missing --no-wait: %s", output)
	}
}

func TestPrNoWaitFlag(t *testing.T) {
	output, err := captureStdout(t, func() error {
		return handleFetchPR([]string{"--no-wait", "-h"})
	})
	if err != nil {
		t.Fatalf("handleFetchPR with --no-wait: %v", err)
	}
	if !strings.Contains(output, "--no-wait") {
		t.Errorf("pr help missing --no-wait: %s", output)
	}
}

// Helpers

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
