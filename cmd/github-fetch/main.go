package main

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/xhd2015/skills/skillcmd"
)

//go:embed SKILL.md
var skillTemplate string

const skillName = "github-fetch"

const help = `
Usage: github-fetch <command> [ARGS]

Commands:
  pr [--diff] [--logs] [--workflow <name>] <url-or-number>  Fetch and display PR content
  ci [--logs] [--workflow <name>] [--run-id <id>] [--job <name>] <url>  Show CI workflow runs and logs
  yaml validate <path>                  Validate a GitHub Actions workflow file
  work-on <url-or-number>              Create a git worktree for the PR
  push [<url-or-number>] [-f]          Push current HEAD to the PR's source branch
  status                               Show GitHub auth status and API rate limits
  skill --show [--header]              Show the content of SKILL.md
  skill --install [<dir>]              Install skill SKILL.md to a directory
  skill --list                         Print the skill name

<url> can be a PR URL, Actions run URL, or Actions job URL.

When <url-or-number> is just a number (e.g. 379), the tool auto-detects
the current git repository's origin URL to construct the full PR URL.

Examples:
  github-fetch pr https://github.com/xhd2015/xgo/pull/379
  github-fetch pr --diff 379
  github-fetch ci 379
  github-fetch ci --logs 379
  github-fetch ci --logs https://github.com/xhd2015/xgo/actions/runs/26795086426/job/78989577716
  github-fetch pr --logs https://github.com/xhd2015/xgo/actions/runs/26795086426/job/78989577716
  github-fetch work-on 379
  github-fetch push -f
  github-fetch skill --install --cursor

Options:
  -h, --help    Show this help message
`

func main() {
	if err := handle(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func handle(args []string) error {
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		fmt.Print(help)
		return nil
	}

	switch args[0] {
	case "pr", "fetch":
		return handleFetchPR(args[1:])
	case "prs", "pulls":
		return handleListPRs(args[1:])
	case "issues":
		return handleListIssues(args[1:])
	case "issue":
		return handleIssue(args[1:])
	case "ci", "checks":
		return handleCI(args[1:])
	case "yaml":
		return handleYAML(args[1:])
	case "work-on":
		return handleWorkon(args[1:])
	case "push":
		return handlePush(args[1:])
	case "status":
		return handleStatus(args[1:])
	case "skill":
		return handleSkill(args[1:])
	default:
		return fmt.Errorf("unknown command: %s", args[0])
	}
}

func handleSkill(args []string) error {
	return (&skillcmd.SingleSkill{
		Name:        skillName,
		RootContent: skillTemplate,
		Usage:       "github-fetch skill --install",
	}).Handle(args)
}
