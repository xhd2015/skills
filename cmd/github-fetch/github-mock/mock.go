package githubmock

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
)

// Server wraps an httptest.Server that mimics GitHub's API.
type Server struct {
	srv  *httptest.Server

	// Configurable responses
	RepoInfo       *RepoInfo
	WorkflowRuns   []WorkflowRun
	WorkflowJobs   []WorkflowJob
	JobLogs        string // plain text, will be zipped
	WorkflowFiles  []string
	DefaultBranch  string // returned in RepoInfo if set

	// PR-related
	PRInfo  *PRInfo
	PRFiles []PRFile
	PRDiff  string

	// Request tracking
	Requests []string
}

// RepoInfo mirrors GitHub's /repos/{owner}/{repo} response subset.
type RepoInfo struct {
	DefaultBranch string `json:"default_branch"`
	FullName      string `json:"full_name"`
}

// WorkflowRun mirrors the subset used by the main package.
type WorkflowRun struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	Conclusion string `json:"conclusion"`
	HTMLURL    string `json:"html_url"`
	HeadBranch string `json:"head_branch"`
	HeadSHA    string `json:"head_sha"`
	CreatedAt  string `json:"created_at"`
	Event      string `json:"event"`
}

// WorkflowJob mirrors the subset used by the main package.
type WorkflowJob struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Status      string `json:"status"`
	Conclusion  string `json:"conclusion"`
	StartedAt   string `json:"started_at"`
	CompletedAt string `json:"completed_at"`
}

// PRInfo mirrors the subset for PR responses.
type PRInfo struct {
	Number  int    `json:"number"`
	Title   string `json:"title"`
	Body    string `json:"body"`
	State   string `json:"state"`
	HTMLURL string `json:"html_url"`
	Head    PRBranch `json:"head"`
	Base    PRBranch `json:"base"`
	User    PRUser   `json:"user"`
}

// PRBranch for mock PR responses.
type PRBranch struct {
	Ref  string `json:"ref"`
	SHA  string `json:"sha"`
	Repo PRRepo `json:"repo"`
}

// PRRepo for mock PR responses.
type PRRepo struct {
	FullName string `json:"full_name"`
	CloneURL string `json:"clone_url"`
	SSHURL   string `json:"ssh_url"`
}

// PRUser for mock PR responses.
type PRUser struct {
	Login string `json:"login"`
}

// PRFile for mock PR file responses.
type PRFile struct {
	Filename  string `json:"filename"`
	Status    string `json:"status"`
	Additions int    `json:"additions"`
	Deletions int    `json:"deletions"`
	Changes   int    `json:"changes"`
	Patch     string `json:"patch"`
}

// NewServer creates and starts a mock GitHub API server.
func NewServer() *Server {
	m := &Server{}
	m.srv = httptest.NewServer(http.HandlerFunc(m.handler))
	return m
}

// URL returns the mock server URL.
func (m *Server) URL() string {
	return m.srv.URL
}

// Close shuts down the mock server.
func (m *Server) Close() {
	m.srv.Close()
}

func (m *Server) handler(w http.ResponseWriter, r *http.Request) {
	m.Requests = append(m.Requests, r.Method+" "+r.URL.Path+"?"+r.URL.RawQuery)

	path := strings.TrimPrefix(r.URL.Path, "/")
	parts := splitPath(path)

	switch {
	// GET /repos/{owner}/{repo}
	case matchPath(path, 3) && r.Method == "GET":
		m.handleRepoInfo(w, r)

	// GET /repos/{owner}/{repo}/pulls/{number}
	case matchPath(path, 4) && indexPart(parts, 3) == "pulls" && r.Method == "GET":
		m.handlePRInfo(w, r)

	// GET /repos/{owner}/{repo}/pulls/{number}/files
	case matchPath(path, 5) && indexPart(parts, 3) == "pulls" && indexPart(parts, 5) == "files" && r.Method == "GET":
		m.handlePRFiles(w, r)

	// GET /repos/{owner}/{repo}/actions/runs
	case strings.Contains(path, "/actions/runs") && !strings.Contains(path, "/jobs") && !strings.Contains(path, "/logs"):
		parts := strings.Split(path, "/")
		if len(parts) == 5 && parts[3] == "actions" && parts[4] == "runs" {
			// /repos/{owner}/{repo}/actions/runs
			m.handleWorkflowRuns(w, r)
		} else if len(parts) == 6 && parts[3] == "actions" && parts[4] == "runs" {
			// /repos/{owner}/{repo}/actions/runs/{id}
			m.handleWorkflowRunDetail(w, r)
		}

	// GET /repos/{owner}/{repo}/actions/runs/{id}/jobs
	case strings.Contains(path, "/actions/runs/") && strings.Contains(path, "/jobs") && r.Method == "GET":
		m.handleWorkflowJobs(w, r)

	// GET /repos/{owner}/{repo}/actions/jobs/{id}/logs
	case strings.Contains(path, "/actions/jobs/") && strings.Contains(path, "/logs") && r.Method == "GET":
		m.handleJobLogs(w, r)

	// GET /repos/{owner}/{repo}/contents/.github/workflows
	case strings.Contains(path, "/contents/.github/workflows") && r.Method == "GET":
		m.handleWorkflowFiles(w, r)

	default:
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "Not Found"})
	}
}

func splitPath(path string) []string {
	return strings.Split(path, "/")
}

func matchPath(path string, n int) bool {
	return len(strings.Split(path, "/")) == n
}

func indexPart(parts []string, i int) string {
	if i < len(parts) {
		return parts[i]
	}
	return ""
}

func (m *Server) handleRepoInfo(w http.ResponseWriter, r *http.Request) {
	info := m.RepoInfo
	if info == nil {
		info = &RepoInfo{DefaultBranch: "main", FullName: "testowner/testrepo"}
	}
	if m.DefaultBranch != "" {
		info.DefaultBranch = m.DefaultBranch
	}
	writeJSON(w, info)
}

func (m *Server) handlePRInfo(w http.ResponseWriter, r *http.Request) {
	if m.PRInfo == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "Not Found"})
		return
	}
	writeJSON(w, m.PRInfo)
}

func (m *Server) handlePRFiles(w http.ResponseWriter, r *http.Request) {
	if m.PRFiles == nil {
		m.PRFiles = []PRFile{}
	}
	writeJSON(w, m.PRFiles)
}

func (m *Server) handleWorkflowRuns(w http.ResponseWriter, r *http.Request) {
	runs := make([]WorkflowRun, len(m.WorkflowRuns))
	copy(runs, m.WorkflowRuns)
	sort.Slice(runs, func(i, j int) bool {
		return runs[i].ID > runs[j].ID
	})
	resp := struct {
		TotalCount   int            `json:"total_count"`
		WorkflowRuns []WorkflowRun `json:"workflow_runs"`
	}{
		TotalCount:   len(runs),
		WorkflowRuns: runs,
	}
	if resp.WorkflowRuns == nil {
		resp.WorkflowRuns = []WorkflowRun{}
	}
	writeJSON(w, resp)
}

func (m *Server) handleWorkflowRunDetail(w http.ResponseWriter, r *http.Request) {
	// Find the run by ID in the path
	path := strings.TrimPrefix(r.URL.Path, "/")
	parts := splitPath(path)
	if len(parts) < 6 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	runID := parts[5]
	for _, run := range m.WorkflowRuns {
		if fmt.Sprintf("%d", run.ID) == runID {
			writeJSON(w, run)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"message": "Not Found"})
}

func (m *Server) handleWorkflowJobs(w http.ResponseWriter, r *http.Request) {
	resp := struct {
		TotalCount int            `json:"total_count"`
		Jobs       []WorkflowJob `json:"jobs"`
	}{
		TotalCount: len(m.WorkflowJobs),
		Jobs:       m.WorkflowJobs,
	}
	if resp.Jobs == nil {
		resp.Jobs = []WorkflowJob{}
	}
	writeJSON(w, resp)
}

func (m *Server) handleJobLogs(w http.ResponseWriter, r *http.Request) {
	logContent := m.JobLogs
	if logContent == "" {
		logContent = "(no logs)"
	}

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	fw, err := zw.Create("log.txt")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fw.Write([]byte(logContent))
	zw.Close()

	w.Header().Set("Content-Type", "application/zip")
	w.Write(buf.Bytes())
}

func (m *Server) handleWorkflowFiles(w http.ResponseWriter, r *http.Request) {
	if m.WorkflowFiles == nil {
		m.WorkflowFiles = []string{}
	}
	files := make([]struct {
		Name string `json:"name"`
		Path string `json:"path"`
	}, len(m.WorkflowFiles))
	for i, f := range m.WorkflowFiles {
		files[i] = struct {
			Name string `json:"name"`
			Path string `json:"path"`
		}{Name: f, Path: ".github/workflows/" + f}
	}
	writeJSON(w, files)
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}
