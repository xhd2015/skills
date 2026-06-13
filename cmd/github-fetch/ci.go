package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/xhd2015/less-flags"
)

const defaultLogLines = 4096

func handleCI(args []string) error {
	var showLogs, fullLogs, noWait bool
	var runID *int64
	var jobFilter, workflowFilter string

	remain, err := lessflags.Bool("--logs", &showLogs).
		Bool("--full", &fullLogs).
		Bool("--no-wait", &noWait).
		String("--job", &jobFilter).
		String("--workflow", &workflowFilter).
		Int("--run-id", &runID).
		Help("-h,--help", ciHelp).
		HelpNoExit().
		Parse(args)
	if err != nil {
		if err == lessflags.ErrHelp {
			return nil
		}
		return err
	}

	if len(remain) == 0 {
		return handleCIAutoDetect(os.Stdout, showLogs, fullLogs, noWait, runID, jobFilter, workflowFilter)
	}
	prRef := remain[0]

	if isActionsURL(prRef) {
		return handleCIActionsURL(os.Stdout, prRef, jobFilter, fullLogs, noWait)
	}

	owner, repo, number, err := resolvePRRef(prRef)
	if err != nil {
		return err
	}

	info, err := fetchPRInfo(owner, repo, number)
	if err != nil {
		return fmt.Errorf("fetch PR #%s: %w", number, err)
	}

	runs, err := fetchWorkflowRuns(owner, repo, info.Head.Ref)
	if err != nil {
		return fmt.Errorf("fetch workflow runs: %w", err)
	}

	if len(runs) == 0 {
		msg := "no workflow runs found for this PR"
		if wfFiles, err := fetchWorkflowFiles(owner, repo); err == nil && len(wfFiles) > 0 {
			quoted := make([]string, len(wfFiles))
			for i, f := range wfFiles {
				quoted[i] = ".github/workflows/" + f
			}
			msg += fmt.Sprintf("\nWorkflow files in repo: %s", strings.Join(quoted, ", "))
			msg += fmt.Sprintf("\n(check for syntax errors: github-fetch yaml validate .github/workflows/%s)", wfFiles[0])
		}
		return fmt.Errorf("%s", msg)
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
		var rid int64
		if runID != nil {
			rid = *runID
		}
		return showRunLogs(os.Stdout, owner, repo, runs, rid, jobFilter, fullLogs, noWait)
	}

	printWorkflowRuns(os.Stdout, info, runs)
	return nil
}

func handleCIAutoDetect(w io.Writer, showLogs, fullLogs, noWait bool, runID *int64, jobFilter, workflowFilter string) error {
	owner, repo, err := getOriginRepo()
	if err != nil {
		return fmt.Errorf("cannot auto-detect repository: %w (use ci <url> instead)", err)
	}
	fmt.Fprintf(w, "Detected repo: https://github.com/%s/%s\n", owner, repo)

	defaultBranch, err := fetchDefaultBranch(owner, repo)
	if err != nil {
		return fmt.Errorf("fetch default branch: %w", err)
	}

	runs, err := fetchRepoWorkflowRuns(owner, repo, defaultBranch)
	if err != nil {
		return fmt.Errorf("fetch workflow runs: %w", err)
	}

	if len(runs) == 0 {
		msg := "no workflow runs found for this repo"
		if wfFiles, err := fetchWorkflowFiles(owner, repo); err == nil && len(wfFiles) > 0 {
			quoted := make([]string, len(wfFiles))
			for i, f := range wfFiles {
				quoted[i] = ".github/workflows/" + f
			}
			msg += fmt.Sprintf("\nWorkflow files in repo: %s", strings.Join(quoted, ", "))
			msg += fmt.Sprintf("\n(check for syntax errors: github-fetch yaml validate .github/workflows/%s)", wfFiles[0])
		}
		return fmt.Errorf("%s", msg)
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
		var rid int64
		if runID != nil {
			rid = *runID
		}
		return showRunLogs(w, owner, repo, runs, rid, jobFilter, fullLogs, noWait)
	}

	printRunsList(w, runs)
	return nil
}

func printRunsList(w io.Writer, runs []WorkflowRun) {
	var bw func(string)
	if f, ok := w.(interface{ Write([]byte) (int, error) }); ok {
		bw = func(s string) { f.Write([]byte(s)) }
	} else {
		sw := w.(interface{ WriteString(string) (int, error) })
		bw = func(s string) { sw.WriteString(s) }
	}

	bw(fmt.Sprintf("Workflow Runs:\n\n"))

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

func waitForRunCompletion(w io.Writer, owner, repo string, runID int64, runName string, noWait bool) error {
	if noWait {
		return nil
	}

	run, err := fetchSingleWorkflowRun(owner, repo, runID)
	if err != nil {
		return fmt.Errorf("check run status: %w", err)
	}

	if run.Status == "completed" {
		return nil
	}

	fmt.Fprintf(w, "Waiting for %s (run #%d) to complete\n", runName, runID)
	fmt.Fprintf(w, "%s %s\n", time.Now().Format("15:04:05"), run.Status)
	prevStatus := run.Status
	startTime := time.Now()

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	timeout := time.After(pollTimeout)

	for {
		select {
		case <-ticker.C:
			run, err = fetchSingleWorkflowRun(owner, repo, runID)
			if err != nil {
				return fmt.Errorf("poll run status: %w", err)
			}
			if run.Status == "completed" {
				fmt.Fprintf(w, "\n%s %s %s\n", time.Now().Format("15:04:05"), run.Status, run.Conclusion)
				elapsed := time.Since(startTime).Round(time.Second)
				fmt.Fprintf(w, "done (took %v)\n", elapsed)
				return nil
			}
			if run.Status != prevStatus {
				fmt.Fprintf(w, "\n%s %s\n", time.Now().Format("15:04:05"), run.Status)
				prevStatus = run.Status
			} else {
				fmt.Fprint(w, ".")
			}
		case <-timeout:
			fmt.Fprint(w, "\n")
			return fmt.Errorf("timed out waiting for run #%d to complete after %v", runID, pollTimeout)
		}
	}
}

func handleCIActionsURL(w io.Writer, url, jobFilter string, fullLogs, noWait bool) error {
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
		return fetchAndPrintJobLog(w, owner, repo, runID, jobID, fullLogs, noWait)
	}

	runs := []WorkflowRun{{ID: runID, Name: fmt.Sprintf("run #%d", runID)}}
	return showRunLogs(w, owner, repo, runs, runID, jobFilter, fullLogs, noWait)
}

func showRunLogs(w io.Writer, owner, repo string, runs []WorkflowRun, runID int64, jobFilter string, fullLogs, noWait bool) error {
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

	if err := waitForRunCompletion(w, owner, repo, targetRun.ID, targetRun.Name, noWait); err != nil {
		return err
	}

	jobs, err := fetchWorkflowRunJobs(owner, repo, targetRun.ID)
	if err != nil {
		return fmt.Errorf("fetch jobs for run %d: %w", targetRun.ID, err)
	}

	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("Workflow: %s (Run #%d, event: %s, branch: %s) — %s\n", targetRun.Name, targetRun.ID, targetRun.Event, targetRun.HeadBranch, targetRun.Conclusion))
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

func fetchAndPrintJobLog(w io.Writer, owner, repo string, runID, jobID int64, fullLogs, noWait bool) error {
	if err := waitForRunCompletion(w, owner, repo, runID, fmt.Sprintf("run #%d", runID), noWait); err != nil {
		return err
	}

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
Usage: github-fetch ci [--logs] [--full] [--no-wait] [--workflow <name>] [--run-id <id>] [--job <name>] <url>

Show CI workflow runs and logs.

<url> can be:
  - PR URL:         https://github.com/owner/repo/pull/123
  - Run URL:        https://github.com/owner/repo/actions/runs/456
  - Job URL:        https://github.com/owner/repo/actions/runs/456/job/789

Options:
  --logs            Show logs (default for job URLs)
  --full            Show complete log output (default: last 4096 lines)
  --no-wait         Skip waiting for in-progress runs to complete
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
