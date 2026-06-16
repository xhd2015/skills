package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/xhd2015/less-flags"
)

func handleFetchPR(args []string) error {
	var showDiff, showLogs, fullLogs, noWait, showList bool
	var workflowFilter string

	remain, err := lessflags.Bool("--list", &showList).
		Bool("--diff", &showDiff).
		Bool("--logs", &showLogs).
		Bool("--full", &fullLogs).
		Bool("--no-wait", &noWait).
		String("--workflow", &workflowFilter).
		Help("-h,--help", prHelp).
		HelpNoExit().
		Parse(args)
	if err != nil {
		if err == lessflags.ErrHelp {
			return nil
		}
		return err
	}

	if showList {
		return handleListPRs(remain)
	}

	if len(remain) == 0 {
		return fmt.Errorf("pr requires a PR URL or number")
	}
	prRef := remain[0]

	if isActionsURL(prRef) || showLogs {
		ciArgs := buildCIArgs(prRef, showLogs, workflowFilter, "", fullLogs, noWait)
		return handleCI(ciArgs)
	}

	owner, repo, number, err := resolvePRRef(prRef)
	if err != nil {
		return err
	}

	info, err := fetchPRInfo(owner, repo, number)
	if err != nil {
		return fmt.Errorf("fetch PR #%s: %w", number, err)
	}

	files, err := fetchPRFiles(owner, repo, number)
	if err != nil {
		return fmt.Errorf("fetch PR files: %w", err)
	}

	var diff string
	if showDiff {
		diff, err = fetchPRDiff(owner, repo, number)
		if err != nil {
			return fmt.Errorf("fetch PR diff: %w", err)
		}
	}

	printPR(os.Stdout, info, files, diff, showDiff)
	return nil
}

func buildCIArgs(prRef string, showLogs bool, workflowFilter, jobFilter string, fullLogs, noWait bool) []string {
	var args []string
	if showLogs {
		args = append(args, "--logs")
	}
	if fullLogs {
		args = append(args, "--full")
	}
	if noWait {
		args = append(args, "--no-wait")
	}
	if workflowFilter != "" {
		args = append(args, "--workflow", workflowFilter)
	}
	if jobFilter != "" {
		args = append(args, "--job", jobFilter)
	}
	args = append(args, prRef)
	return args
}

const prHelp = `
Usage: github-fetch pr [--diff] [--logs] [--full] [--no-wait] [--workflow <name>] <url-or-number>

Fetch and display PR content.

Options:
  --diff            Show full diff (default: simplified diff, up to 10 files, 3 lines each)
  --logs            Show CI workflow logs for the PR
  --full            Show complete log output (use with --logs; default: last 4096 lines)
  --no-wait         Skip waiting for in-progress runs to complete (use with --logs)
  --workflow <name> Filter workflow runs by name (use with --logs)
  -h, --help        Show this help message

When <url> is an Actions URL (runs/.../job/...), --logs delegates to the ci subcommand.
`

func printPR(out interface{ WriteString(string) (int, error) }, info *PRInfo, files []PRFile, diff string, showDiff bool) {
	var w func(string)
	if f, ok := out.(interface{ Write([]byte) (int, error) }); ok {
		w = func(s string) { f.Write([]byte(s)) }
	} else {
		sw := out.(interface{ WriteString(string) (int, error) })
		w = func(s string) { sw.WriteString(s) }
	}

	stateLabel := info.State
	switch info.State {
	case "open":
		stateLabel = "\033[32mopen\033[0m"
	case "closed":
		stateLabel = "\033[31mclosed\033[0m"
	case "merged":
		stateLabel = "\033[35mmerged\033[0m"
	}

	w(fmt.Sprintf("PR #%d: %s\n", info.Number, info.Title))
	w(fmt.Sprintf("State:  %s\n", stateLabel))
	w(fmt.Sprintf("Author: @%s\n", info.User.Login))
	w(fmt.Sprintf("Branch: %s → %s:%s\n", info.Head.Repo.FullName, info.Base.Repo.FullName, info.Base.Ref))
	w(fmt.Sprintf("URL:    %s\n", info.HTMLURL))

	if info.Body != "" {
		w("\n── Description ──\n")
		w(info.Body)
		if !strings.HasSuffix(info.Body, "\n") {
			w("\n")
		}
	}

	if len(files) > 0 {
		totalAdd, totalDel := 0, 0
		for _, f := range files {
			totalAdd += f.Additions
			totalDel += f.Deletions
		}

		w("\n── Files Changed ──\n")
		for _, f := range files {
			w(fmt.Sprintf("  %s  %s  (+%d -%d)\n", f.Status, f.Filename, f.Additions, f.Deletions))
		}
		w("  ─────────────────\n")
		w(fmt.Sprintf("  %d files changed, +%d -%d\n", len(files), totalAdd, totalDel))
	}

	if showDiff {
		if diff != "" {
			w("\n── Diff ──\n")
			w(diff)
			if !strings.HasSuffix(diff, "\n") {
				w("\n")
			}
		}
	} else {
		printSimplifiedDiff(w, files)
	}
}

func printSimplifiedDiff(w func(string), files []PRFile) {
	maxFiles := 10
	if len(files) < maxFiles {
		maxFiles = len(files)
	}
	if maxFiles == 0 {
		return
	}

	w("\n── Diff (simplified) ──\n")
	shown := 0
	for _, f := range files {
		if shown >= maxFiles {
			break
		}
		if f.Patch == "" {
			continue
		}
		lines := strings.Split(strings.TrimSuffix(f.Patch, "\n"), "\n")
		if len(lines) == 0 {
			continue
		}
		shown++
		w(fmt.Sprintf("  %s:\n", f.Filename))
		limit := 3
		if len(lines) < limit {
			limit = len(lines)
		}
		for i := 0; i < limit; i++ {
			w(fmt.Sprintf("    %s\n", lines[i]))
		}
		if len(lines) > 3 {
			w("    ...\n")
		}
	}
	if len(files) > 10 {
		w(fmt.Sprintf("  ... and %d more files\n", len(files)-10))
	}
}
