package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/xhd2015/less-flags"
)

func handleListIssues(args []string) error {
	opts, remain, err := parseListFlags(args)
	if err != nil {
		return err
	}

	owner, repo, err := resolveListRepo(remain)
	if err != nil {
		return err
	}

	items, hasNext, err := fetchIssueList(owner, repo, opts.state, opts.page, opts.perPage)
	if err != nil {
		return fmt.Errorf("fetch issues: %w", err)
	}

	printIssueList(os.Stdout, owner, repo, opts.state, opts.page, items, hasNext)
	return nil
}

func handleIssue(args []string) error {
	var showList bool
	remain, err := lessflags.Bool("--list", &showList).Parse(args)
	if err != nil {
		return err
	}
	if showList {
		return handleListIssues(remain)
	}
	return fmt.Errorf("issue requires --list (or use `issues`)")
}

func fetchIssueList(owner, repo, state string, page, perPage int) ([]IssueListItem, bool, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/issues?state=%s&page=%d&per_page=%d",
		getAPIBaseURL(), owner, repo, state, page, perPage)

	data, headers, err := apiGetWithHeaders(url)
	if err != nil {
		return nil, false, err
	}

	var resp []githubIssueListResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, false, fmt.Errorf("parse issue list: %w", err)
	}

	var items []IssueListItem
	for _, issue := range resp {
		if issue.PullRequest != nil {
			continue
		}
		labels := make([]string, len(issue.Labels))
		for i, label := range issue.Labels {
			labels[i] = label.Name
		}
		items = append(items, IssueListItem{
			Number:  issue.Number,
			Title:   issue.Title,
			State:   issue.State,
			User:    issue.User.Login,
			HTMLURL: issue.HTMLURL,
			Labels:  labels,
		})
	}

	return items, hasNextPage(headers.Get("Link")), nil
}

func printIssueList(out interface{ Write([]byte) (int, error) }, owner, repo, state string, page int, items []IssueListItem, hasNext bool) {
	w := func(s string) { out.Write([]byte(s)) }

	if len(items) == 0 {
		w(emptyIssueMessage(state))
		return
	}

	printListHeader(w, owner, repo, state)
	for _, issue := range items {
		printListItem(w, issue.Number, issue.State, issue.User, issue.Title)
	}
	printPaginationFooter(w, page, len(items), hasNext)
}

func emptyIssueMessage(state string) string {
	switch state {
	case "open":
		return "no open issues\n"
	case "closed":
		return "no closed issues\n"
	default:
		return "no issues\n"
	}
}