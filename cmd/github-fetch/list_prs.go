package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func handleListPRs(args []string) error {
	opts, remain, err := parseListFlags(args)
	if err != nil {
		return err
	}

	owner, repo, err := resolveListRepo(remain)
	if err != nil {
		return err
	}

	items, hasNext, err := fetchPRList(owner, repo, opts.state, opts.page, opts.perPage)
	if err != nil {
		return fmt.Errorf("fetch pull requests: %w", err)
	}

	printPRList(os.Stdout, owner, repo, opts.state, opts.page, items, hasNext)
	return nil
}

func fetchPRList(owner, repo, state string, page, perPage int) ([]PRListItem, bool, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/pulls?state=%s&page=%d&per_page=%d",
		getAPIBaseURL(), owner, repo, state, page, perPage)

	data, headers, err := apiGetWithHeaders(url)
	if err != nil {
		return nil, false, err
	}

	var resp []githubPRListResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, false, fmt.Errorf("parse pull request list: %w", err)
	}

	items := make([]PRListItem, len(resp))
	for i, pr := range resp {
		items[i] = PRListItem{
			Number:  pr.Number,
			Title:   pr.Title,
			State:   pr.State,
			User:    pr.User.Login,
			HTMLURL: pr.HTMLURL,
		}
	}

	return items, hasNextPage(headers.Get("Link")), nil
}

func printPRList(out interface{ Write([]byte) (int, error) }, owner, repo, state string, page int, items []PRListItem, hasNext bool) {
	w := func(s string) { out.Write([]byte(s)) }

	if len(items) == 0 {
		w(emptyPRMessage(state))
		return
	}

	printListHeader(w, owner, repo, state)
	for _, pr := range items {
		printListItem(w, pr.Number, pr.State, pr.User, pr.Title)
	}
	printPaginationFooter(w, page, len(items), hasNext)
}

func emptyPRMessage(state string) string {
	switch state {
	case "open":
		return "no open pull requests\n"
	case "closed":
		return "no closed pull requests\n"
	default:
		return "no pull requests\n"
	}
}