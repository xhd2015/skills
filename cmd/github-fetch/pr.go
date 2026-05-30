package main

import (
	"fmt"
	"os"
	"strings"
)

func handleFetchPR(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("pr requires a PR URL or number")
	}
	owner, repo, number, err := resolvePRRef(args[0])
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

	diff, err := fetchPRDiff(owner, repo, number)
	if err != nil {
		return fmt.Errorf("fetch PR diff: %w", err)
	}

	printPR(os.Stdout, info, files, diff)
	return nil
}

func printPR(out interface{ WriteString(string) (int, error) }, info *PRInfo, files []PRFile, diff string) {
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
		w("\n── Files Changed ──\n")
		totalAdd, totalDel := 0, 0
		for _, f := range files {
			totalAdd += f.Additions
			totalDel += f.Deletions
		}
		for _, f := range files {
			w(fmt.Sprintf("  %s  %s  (+%d -%d)\n", f.Status, f.Filename, f.Additions, f.Deletions))
		}
		w(fmt.Sprintf("  ─────────────────\n"))
		w(fmt.Sprintf("  %d files changed, +%d -%d\n", len(files), totalAdd, totalDel))
	}

	if diff != "" {
		w("\n── Diff ──\n")
		w(diff)
		if !strings.HasSuffix(diff, "\n") {
			w("\n")
		}
	}
}
