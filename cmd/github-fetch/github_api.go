package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type PRInfo struct {
	Number  int
	Title   string
	Body    string
	State   string
	HTMLURL string
	Head    PRBranch
	Base    PRBranch
	User    PRUser
}

type PRBranch struct {
	Ref  string
	SHA  string
	Repo PRRepo
}

type PRRepo struct {
	FullName string
	CloneURL string
	SSHURL   string
}

type PRUser struct {
	Login string
}

type PRFile struct {
	Filename  string
	Status    string
	Additions int
	Deletions int
	Changes   int
	Patch     string
}

type githubPRResponse struct {
	Number    int              `json:"number"`
	Title     string           `json:"title"`
	Body      string           `json:"body"`
	State     string           `json:"state"`
	HTMLURL   string           `json:"html_url"`
	Head      githubPRBranch   `json:"head"`
	Base      githubPRBranch   `json:"base"`
	User      githubPRUser     `json:"user"`
}

type githubPRBranch struct {
	Ref  string       `json:"ref"`
	SHA  string       `json:"sha"`
	Repo githubPRRepo `json:"repo"`
}

type githubPRRepo struct {
	FullName string `json:"full_name"`
	CloneURL string `json:"clone_url"`
	SSHURL   string `json:"ssh_url"`
}

type githubPRUser struct {
	Login string `json:"login"`
}

type githubPRFile struct {
	Filename  string `json:"filename"`
	Status    string `json:"status"`
	Additions int    `json:"additions"`
	Deletions int    `json:"deletions"`
	Changes   int    `json:"changes"`
	Patch     string `json:"patch"`
}

var (
	githubURLRe  = regexp.MustCompile(`^https?://github\.com/([^/]+)/([^/]+)/pull/(\d+)(?:/.*)?$`)
	gitSSHRe     = regexp.MustCompile(`^git@github\.com:([^/]+)/([^/.]+?)(?:\.git)?$`)
	gitSSHURLRe  = regexp.MustCompile(`^ssh://git@github\.com/([^/]+)/([^/.]+?)(?:\.git)?$`)
	gitHTTPSRe   = regexp.MustCompile(`^https?://github\.com/([^/]+)/([^/.]+?)(?:\.git)?$`)
)

func resolvePRRef(raw string) (owner, repo, number string, err error) {
	if n, convErr := strconv.Atoi(raw); convErr == nil {
		owner, repo, err = getOriginRepo()
		if err != nil {
			return "", "", "", fmt.Errorf("cannot auto-detect PR from number: %w", err)
		}
		return owner, repo, strconv.Itoa(n), nil
	}
	return parseGitHubURL(raw)
}

func parseGitHubURL(raw string) (owner, repo, number string, err error) {
	m := githubURLRe.FindStringSubmatch(raw)
	if m == nil {
		return "", "", "", fmt.Errorf("invalid PR reference: %q (expected URL like https://github.com/owner/repo/pull/123 or a PR number)", raw)
	}
	return m[1], m[2], m[3], nil
}

func getOriginRepo() (owner, repo string, err error) {
	out, err := runGitCmd("remote", "get-url", "origin")
	if err != nil {
		return "", "", fmt.Errorf("get origin URL: %w", err)
	}
	return parseOriginURL(strings.TrimSpace(out))
}

func parseOriginURL(raw string) (owner, repo string, err error) {
	if m := gitSSHRe.FindStringSubmatch(raw); m != nil {
		return m[1], m[2], nil
	}
	if m := gitSSHURLRe.FindStringSubmatch(raw); m != nil {
		return m[1], m[2], nil
	}
	if m := gitHTTPSRe.FindStringSubmatch(raw); m != nil {
		return m[1], m[2], nil
	}
	return "", "", fmt.Errorf("cannot parse origin URL: %s", raw)
}

func fetchPRInfo(owner, repo, number string) (*PRInfo, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls/%s", owner, repo, number)
	data, err := apiGet(url)
	if err != nil {
		return nil, err
	}
	var resp githubPRResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parse PR response: %w", err)
	}
	return &PRInfo{
		Number:  resp.Number,
		Title:   resp.Title,
		Body:    resp.Body,
		State:   resp.State,
		HTMLURL: resp.HTMLURL,
		Head: PRBranch{
			Ref: resp.Head.Ref,
			SHA: resp.Head.SHA,
			Repo: PRRepo{
				FullName: resp.Head.Repo.FullName,
				CloneURL: resp.Head.Repo.CloneURL,
				SSHURL:   resp.Head.Repo.SSHURL,
			},
		},
		Base: PRBranch{
			Ref: resp.Base.Ref,
			SHA: resp.Base.SHA,
			Repo: PRRepo{
				FullName: resp.Base.Repo.FullName,
				CloneURL: resp.Base.Repo.CloneURL,
				SSHURL:   resp.Base.Repo.SSHURL,
			},
		},
		User: PRUser{
			Login: resp.User.Login,
		},
	}, nil
}

func fetchPRFiles(owner, repo, number string) ([]PRFile, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls/%s/files", owner, repo, number)
	data, err := apiGet(url)
	if err != nil {
		return nil, err
	}
	var files []githubPRFile
	if err := json.Unmarshal(data, &files); err != nil {
		return nil, fmt.Errorf("parse PR files response: %w", err)
	}
	result := make([]PRFile, len(files))
	for i, f := range files {
		result[i] = PRFile{
			Filename:  f.Filename,
			Status:    f.Status,
			Additions: f.Additions,
			Deletions: f.Deletions,
			Changes:   f.Changes,
			Patch:     f.Patch,
		}
	}
	return result, nil
}

func fetchPRDiff(owner, repo, number string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls/%s", owner, repo, number)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github.v3.diff")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return tryGhDiff(owner, repo, number, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusForbidden {
		return tryGhDiff(owner, repo, number, fmt.Errorf("HTTP %d (repo may be private)", resp.StatusCode))
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read diff response: %w", err)
	}
	return string(data), nil
}

func apiGet(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return tryGhAPI(url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusForbidden {
		return tryGhAPI(url, fmt.Errorf("HTTP %d", resp.StatusCode))
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}
	return data, nil
}

func tryGhAPI(url string, originalErr error) ([]byte, error) {
	if !isGHAvailable() {
		return nil, originalErr
	}
	apiPath := strings.TrimPrefix(url, "https://api.github.com/")
	out, err := runCmd("gh", "api", apiPath)
	if err != nil {
		return nil, fmt.Errorf("gh api failed: %w (original: %v)", err, originalErr)
	}
	return []byte(out), nil
}

func tryGhDiff(owner, repo, number string, originalErr error) (string, error) {
	if !isGHAvailable() {
		return "", originalErr
	}
	apiPath := fmt.Sprintf("repos/%s/%s/pulls/%s", owner, repo, number)
	out, err := runCmd("gh", "api", "-H", "Accept: application/vnd.github.v3.diff", apiPath)
	if err != nil {
		return "", fmt.Errorf("gh api diff failed: %w (original: %v)", err, originalErr)
	}
	return out, nil
}

func isGHAvailable() bool {
	_, err := exec.LookPath("gh")
	return err == nil
}

func getCurrentBranch() (string, error) {
	out, err := runGitCmd("rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", fmt.Errorf("get current branch: %w", err)
	}
	return strings.TrimSpace(out), nil
}

var prBranchRe = regexp.MustCompile(`^pr-(\d+)$`)

func parsePRBranch(branch string) (number string, ok bool) {
	m := prBranchRe.FindStringSubmatch(branch)
	if m == nil {
		return "", false
	}
	return m[1], true
}

func resolvePRFromBranch() (owner, repo, number string, err error) {
	branch, err := getCurrentBranch()
	if err != nil {
		return "", "", "", err
	}
	num, ok := parsePRBranch(branch)
	if !ok {
		return "", "", "", fmt.Errorf("not on a PR branch (expected pr-<number>, got %q), specify PR URL or number", branch)
	}
	owner, repo, err = getOriginRepo()
	if err != nil {
		return "", "", "", fmt.Errorf("cannot auto-detect repo from branch: %w", err)
	}
	return owner, repo, num, nil
}

func runGitCmd(args ...string) (string, error) {
	return runCmd("git", args...)
}

func runCmd(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("%s %s: %s", name, strings.Join(args, " "), strings.TrimSpace(string(exitErr.Stderr)))
		}
		return "", fmt.Errorf("%s %s: %w", name, strings.Join(args, " "), err)
	}
	return string(out), nil
}
