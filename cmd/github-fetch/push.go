package main

import (
	"fmt"
	"strings"
)

func handlePush(args []string) error {
	force := false
	var prRef string

	i := 0
	for i < len(args) {
		if args[i] == "-f" || args[i] == "--force" {
			force = true
			i++
			continue
		}
		if args[i] == "-h" || args[i] == "--help" {
			fmt.Print(pushHelp)
			return nil
		}
		break
	}
	if i < len(args) {
		prRef = args[i]
	}

	var owner, repo, number string
	var err error

	if prRef != "" {
		owner, repo, number, err = resolvePRRef(prRef)
	} else {
		owner, repo, number, err = resolvePRFromBranch()
	}
	if err != nil {
		return err
	}

	info, err := fetchPRInfo(owner, repo, number)
	if err != nil {
		return fmt.Errorf("fetch PR #%s: %w", number, err)
	}

	pushURL := info.Head.Repo.CloneURL
	if pushURL == "" {
		pushURL = info.Head.Repo.SSHURL
	}
	if pushURL == "" {
		return fmt.Errorf("no push URL for PR head repository %s", info.Head.Repo.FullName)
	}

	pushArgs := []string{"push"}
	if force {
		pushArgs = append(pushArgs, "-f")
	}
	pushArgs = append(pushArgs, pushURL, fmt.Sprintf("HEAD:%s", info.Head.Ref))

	fmt.Printf("Pushing to %s:%s...\n", info.Head.Repo.FullName, info.Head.Ref)
	out, err := runGitCmd(pushArgs...)
	if err != nil {
		return fmt.Errorf("push failed: %w", err)
	}
	if out != "" {
		fmt.Print(strings.TrimSpace(out))
		fmt.Println()
	}
	fmt.Printf("Pushed to PR #%d: %s\n", info.Number, info.HTMLURL)
	return nil
}

const pushHelp = `
Usage: github-fetch push [<url-or-number>] [-f]

Push the current HEAD to the PR's source branch.

When <url-or-number> is omitted, the current branch name is used to
auto-detect the PR (must match pr-<number>).

Options:
  -f, --force   Force push (git push -f)
  -h, --help    Show this help message
`
