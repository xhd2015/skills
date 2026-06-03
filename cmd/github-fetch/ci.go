package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const defaultLogLines = 4096

func handleCI(args []string) error {
	showLogs := false
	fullLogs := false
	var runID int64
	var jobFilter string
	var workflowFilter string
	var prRef string

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--logs":
			showLogs = true
		case "--full":
			fullLogs = true
		case "--run-id":
			i++
			if i >= len(args) {
				return fmt.Errorf("--run-id requires a value")
			}
			var err error
			runID, err = strconv.ParseInt(args[i], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid --run-id value: %s", args[i])
			}
		case "--job":
			i++
			if i >= len(args) {
				return fmt.Errorf("--job requires a value")
			}
			jobFilter = args[i]
		case "--workflow":
			i++
			if i >= len(args) {
				return fmt.Errorf("--workflow requires a value")
			}
			workflowFilter = args[i]
		case "-h", "--help":
			fmt.Print(ciHelp)
			return nil
		default:
			if strings.HasPrefix(args[i], "-") {
				return fmt.Errorf("unknown flag: %s", args[i])
			}
			prRef = args[i]
		}
	}

	if prRef == "" {
		return fmt.Errorf("ci requires a PR URL, Actions run URL, or Actions job URL")
	}

	if isActionsURL(prRef) {
		return handleCIActionsURL(os.Stdout, prRef, jobFilter, fullLogs)
	}

	owner, repo, number, err := resolvePRRef(prRef)
	if err != nil {
		return err
	}

	info, err := fetchPRInfo(owner, repo, number)
	if err != nil {
		return fmt.Errorf("fetch PR #%s: %w", number, err)
	}

	runs, err := fetchWorkflowRuns(owner, repo, info.Head.SHA)
	if err != nil {
		return fmt.Errorf("fetch workflow runs: %w", err)
	}

	if workflowFilter != "" {
		filtered := filterRunsByName(runs, workflowFilter)
		if len(filtered) == 0 {
			names := make([]string, len(runs))
			for i, r := range runs {
				names[i] = r.Name
			}
			return fmt.Errorf("no workflow runs matching %q\nAvailable workflows: %s", workflowFilter, strings.Join(names, ", "))
		}
		runs = filtered
	}

	if showLogs {
		return showRunLogs(os.Stdout, owner, repo, runs, runID, jobFilter, fullLogs)
	}

	printWorkflowRuns(os.Stdout, info, runs)
	return nil
}

func filterRunsByName(runs []WorkflowRun, nameFilter string) []WorkflowRun {
	lower := strings.ToLower(normalizeSlashes(nameFilter))
	var filtered []WorkflowRun
	for _, r := range runs {
		if strings.Contains(strings.ToLower(normalizeSlashes(r.Name)), lower) {
			filtered = append(filtered, r)
		}
	}
	return filtered
}

func normalizeSlashes(s string) string {
	parts := strings.Split(s, "/")
	for i, p := range parts {
		parts[i] = strings.TrimSpace(p)
	}
	return strings.Join(parts, "/")
}

func handleCIActionsURL(w io.Writer, url, jobFilter string, fullLogs bool) error {
	owner, repo, runIDStr, jobIDStr, err := resolveActionsURL(url)
	if err != nil {
		return err
	}

	runID, err := strconv.ParseInt(runIDStr, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid run ID: %s", runIDStr)
	}

	if jobIDStr != "" {
		jobID, err := strconv.ParseInt(jobIDStr, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid job ID: %s", jobIDStr)
		}
		return fetchAndPrintJobLog(w, owner, repo, runID, jobID, fullLogs)
	}

	runs := []WorkflowRun{{ID: runID, Name: fmt.Sprintf("run #%d", runID)}}
	return showRunLogs(w, owner, repo, runs, runID, jobFilter, fullLogs)
}

func showRunLogs(w io.Writer, owner, repo string, runs []WorkflowRun, runID int64, jobFilter string, fullLogs bool) error {
	var targetRun *WorkflowRun
	if runID != 0 {
		for i := range runs {
			if runs[i].ID == runID {
				targetRun = &runs[i]
				break
			}
		}
		if targetRun == nil {
			return fmt.Errorf("workflow run %d not found (use `github-fetch ci <pr>` to list runs)", runID)
		}
	} else {
		if len(runs) == 0 {
			return fmt.Errorf("no workflow runs found")
		}
		targetRun = &runs[0]
	}

	jobs, err := fetchWorkflowRunJobs(owner, repo, targetRun.ID)
	if err != nil {
		return fmt.Errorf("fetch jobs for run %d: %w", targetRun.ID, err)
	}

	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("Logs for %s (run #%d):\n", targetRun.Name, targetRun.ID))
	buf.WriteString(fmt.Sprintf("URL: %s\n\n", targetRun.HTMLURL))

	if len(jobs) == 0 {
		buf.WriteString("  No jobs found.\n")
		io.WriteString(w, buf.String())
		return nil
	}

	for _, j := range jobs {
		if jobFilter != "" && !strings.Contains(strings.ToLower(j.Name), strings.ToLower(jobFilter)) {
			continue
		}

		logContent, err := fetchJobLogs(owner, repo, j.ID)
		if err != nil {
			return showRunLogsViaGh(w, buf.String(), owner, repo, targetRun.ID, fullLogs, err)
		}

		buf.WriteString(fmt.Sprintf("Job: %s (%s)\n", j.Name, j.Conclusion))
		buf.WriteString(strings.Repeat("─", 60) + "\n")
		buf.WriteString(truncateLog(logContent, fullLogs))
		buf.WriteString(strings.Repeat("─", 60) + "\n\n")
	}

	io.WriteString(w, buf.String())
	return nil
}

func showRunLogsViaGh(w io.Writer, partial string, owner, repo string, runID int64, fullLogs bool, apiErr error) error {
	ghOutput, err := ghRunViewLogs(owner, repo, runID)
	if err != nil {
		return fmt.Errorf("fetch logs via API: %v; gh run view --log also failed: %v", apiErr, err)
	}
	io.WriteString(w, partial)
	io.WriteString(w, "\n(Falling back to gh run view --log)\n\n")
	io.WriteString(w, truncateLog(ghOutput, fullLogs))
	return nil
}

func fetchAndPrintJobLog(w io.Writer, owner, repo string, runID, jobID int64, fullLogs bool) error {
	logContent, err := fetchJobLogs(owner, repo, jobID)
	if err != nil {
		ghOutput, ghErr := ghRunViewLogs(owner, repo, runID)
		if ghErr != nil {
			return fmt.Errorf("fetch logs for job %d: API returned %v; gh fallback also failed: %v", jobID, err, ghErr)
		}
		io.WriteString(w, fmt.Sprintf("Logs for job #%d (via gh run view --log):\n\n", jobID))
		io.WriteString(w, strings.Repeat("─", 60) + "\n")
		io.WriteString(w, truncateLog(ghOutput, fullLogs))
		io.WriteString(w, strings.Repeat("─", 60) + "\n")
		return nil
	}

	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("Logs for job #%d:\n\n", jobID))
	buf.WriteString(strings.Repeat("─", 60) + "\n")
	buf.WriteString(truncateLog(logContent, fullLogs))
	buf.WriteString(strings.Repeat("─", 60) + "\n")

	io.WriteString(w, buf.String())
	return nil
}

func truncateLog(content string, full bool) string {
	if full || content == "" {
		return content
	}
	lines := strings.Split(content, "\n")
	if len(lines) <= defaultLogLines {
		return content
	}
	truncated := strings.Join(lines[len(lines)-defaultLogLines:], "\n")
	return fmt.Sprintf("(showing last %d of %d lines; use --full for complete output)\n\n%s", defaultLogLines, len(lines), truncated)
}

const ciHelp = `
Usage: github-fetch ci [--logs] [--full] [--workflow <name>] [--run-id <id>] [--job <name>] <url>

Show CI workflow runs and logs.

<url> can be:
  - PR URL:         https://github.com/owner/repo/pull/123
  - Run URL:        https://github.com/owner/repo/actions/runs/456
  - Job URL:        https://github.com/owner/repo/actions/runs/456/job/789

Options:
  --logs            Show logs (default for job URLs)
  --full            Show complete log output (default: last 4096 lines)
  --workflow <name> Filter workflow runs by name (case-insensitive substring match)
  --run-id <id>     Target a specific workflow run (use with PR URLs)
  --job <name>      Filter logs to a specific job name
  -h, --help        Show this help message
`

func printWorkflowRuns(w io.Writer, info *PRInfo, runs []WorkflowRun) {
	var bw func(string)
	if f, ok := w.(interface{ Write([]byte) (int, error) }); ok {
		bw = func(s string) { f.Write([]byte(s)) }
	} else {
		sw := w.(interface{ WriteString(string) (int, error) })
		bw = func(s string) { sw.WriteString(s) }
	}

	bw(fmt.Sprintf("Workflow Runs for PR #%d (%s):\n\n", info.Number, info.Head.Ref))

	if len(runs) == 0 {
		bw("  No workflow runs found.\n")
		return
	}

	for _, r := range runs {
		icon := "○"
		switch r.Conclusion {
		case "success":
			icon = "\033[32m✓\033[0m"
		case "failure":
			icon = "\033[31m✗\033[0m"
		case "cancelled":
			icon = "\033[33m✗\033[0m"
		}
		if r.Status == "in_progress" {
			icon = "\033[33m●\033[0m"
		} else if r.Status == "queued" || r.Status == "waiting" || r.Status == "pending" {
			icon = "\033[36m○\033[0m"
		}

		status := r.Conclusion
		if r.Conclusion == "" || r.Status == "in_progress" || r.Status == "queued" || r.Status == "waiting" || r.Status == "pending" {
			status = r.Status
		}

		bw(fmt.Sprintf("  %s  %s  %s  %s\n", icon, fmt.Sprintf("%-24s", r.Name), fmt.Sprintf("%-12s", status), r.HTMLURL))
	}
}
