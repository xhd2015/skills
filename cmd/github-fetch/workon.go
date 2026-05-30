package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func handleWorkon(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("work-on requires a PR URL or number")
	}
	owner, repo, number, err := resolvePRRef(args[0])
	if err != nil {
		return err
	}

	info, err := fetchPRInfo(owner, repo, number)
	if err != nil {
		return fmt.Errorf("fetch PR #%s: %w", number, err)
	}

	if err := verifyOriginMatches(info.Base.Repo); err != nil {
		return err
	}

	branchName := fmt.Sprintf("pr-%s", number)
	fmt.Printf("Fetching PR #%s from %s/%s...\n", number, info.Head.Repo.FullName, info.Head.Ref)

	if _, err := runGitCmd("fetch", "origin", fmt.Sprintf("pull/%s/head:%s", number, branchName)); err != nil {
		return fmt.Errorf("fetch PR branch: %w", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get current directory: %w", err)
	}
	repoName := filepath.Base(cwd)
	worktreePath := filepath.Join(filepath.Dir(cwd), fmt.Sprintf("%s-pr-%s", repoName, number))

	fmt.Printf("Creating worktree at %s...\n", worktreePath)
	if _, err := runGitCmd("worktree", "add", worktreePath, branchName); err != nil {
		return fmt.Errorf("create worktree: %w", err)
	}

	fmt.Printf("Worktree created: %s\n", worktreePath)
	fmt.Printf("PR #%d: %s (by @%s)\n", info.Number, info.Title, info.User.Login)
	return nil
}

func verifyOriginMatches(baseRepo PRRepo) error {
	originOwner, originRepo, err := getOriginRepo()
	if err != nil {
		return fmt.Errorf("get origin: %w", err)
	}
	originFullName := originOwner + "/" + originRepo
	if !strings.EqualFold(originFullName, baseRepo.FullName) {
		return fmt.Errorf(
			"origin repository %q does not match PR base repository %q",
			originFullName, baseRepo.FullName,
		)
	}
	return nil
}
