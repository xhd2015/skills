package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const statusHelp = `
Usage: github-fetch status

Show GitHub authentication status and API rate limits.

Options:
  -h, --help    Show this help message
`

type AuthStatus struct {
	GhAvailable    bool
	GhLoggedIn     bool
	GhHost         string
	GhUsername     string
	GhTokenScopes  []string
	EnvTokenSet    bool
	EnvTokenMasked string
	APIBaseURL     string
	APIAuthMethod  AuthMethod
	APIAuthenticated bool
	APIUsername    string
	RateLimit      int
	RateRemaining  int
	RateReset      time.Time
}

type githubUserResponse struct {
	Login string `json:"login"`
}

type githubRateLimitResponse struct {
	Resources struct {
		Core struct {
			Limit     int   `json:"limit"`
			Remaining int   `json:"remaining"`
			Reset     int64 `json:"reset"`
		} `json:"core"`
	} `json:"resources"`
}

func handleStatus(args []string) error {
	if len(args) > 0 && (args[0] == "-h" || args[0] == "--help") {
		fmt.Print(statusHelp)
		return nil
	}
	if len(args) > 0 {
		return fmt.Errorf("unknown status argument: %s", args[0])
	}

	status, err := gatherAuthStatus()
	if err != nil {
		return err
	}
	printAuthStatus(os.Stdout, status)
	return nil
}

func gatherAuthStatus() (*AuthStatus, error) {
	status := &AuthStatus{
		APIBaseURL: getAPIBaseURL(),
	}

	ghAvailable, ghLoggedIn, ghHost, ghUsername, ghScopes := getGhAuthInfo()
	status.GhAvailable = ghAvailable
	status.GhLoggedIn = ghLoggedIn
	status.GhHost = ghHost
	status.GhUsername = ghUsername
	status.GhTokenScopes = ghScopes

	if envToken := strings.TrimSpace(os.Getenv("GITHUB_TOKEN")); envToken != "" {
		status.EnvTokenSet = true
		status.EnvTokenMasked = maskToken(envToken)
	}

	token, method := resolveAuthToken()
	status.APIAuthMethod = method
	status.APIAuthenticated = token != ""

	if status.APIAuthenticated {
		user, err := probeGitHubUser()
		if err != nil {
			return nil, fmt.Errorf("API user probe failed: %w", err)
		}
		status.APIUsername = user
	}

	rateLimit, rateRemaining, rateReset, err := probeRateLimit()
	if err != nil {
		return nil, fmt.Errorf("API rate limit probe failed: %w", err)
	}
	status.RateLimit = rateLimit
	status.RateRemaining = rateRemaining
	status.RateReset = rateReset

	return status, nil
}

func probeGitHubUser() (string, error) {
	url := getAPIBaseURL() + "/user"
	resp, err := authenticatedGet(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read user response: %w", err)
	}

	var user githubUserResponse
	if err := json.Unmarshal(data, &user); err != nil {
		return "", fmt.Errorf("parse user response: %w", err)
	}
	return user.Login, nil
}

func probeRateLimit() (limit, remaining int, reset time.Time, err error) {
	url := getAPIBaseURL() + "/rate_limit"
	resp, err := authenticatedGet(url)
	if err != nil {
		return 0, 0, time.Time{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return 0, 0, time.Time{}, fmt.Errorf("HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, time.Time{}, fmt.Errorf("read rate limit response: %w", err)
	}

	var rateResp githubRateLimitResponse
	if err := json.Unmarshal(data, &rateResp); err != nil {
		return 0, 0, time.Time{}, fmt.Errorf("parse rate limit response: %w", err)
	}

	reset = time.Unix(rateResp.Resources.Core.Reset, 0).UTC()
	return rateResp.Resources.Core.Limit, rateResp.Resources.Core.Remaining, reset, nil
}

func printAuthStatus(w io.Writer, status *AuthStatus) {
	fmt.Fprintln(w, "Auth Status")
	fmt.Fprintln(w, "───────────")

	if status.GhAvailable {
		fmt.Fprintln(w, "gh CLI:        available")
		if status.GhLoggedIn {
			if status.GhHost != "" {
				fmt.Fprintf(w, "GitHub host:   %s\n", status.GhHost)
			}
			if status.GhUsername != "" {
				fmt.Fprintf(w, "Logged in as:  %s\n", status.GhUsername)
			}
			if len(status.GhTokenScopes) > 0 {
				fmt.Fprintf(w, "Token scopes:  %s\n", strings.Join(status.GhTokenScopes, ", "))
			}
		}
	} else {
		fmt.Fprintln(w, "gh CLI:        not found")
	}

	if status.EnvTokenSet {
		fmt.Fprintf(w, "GITHUB_TOKEN:  set (%s)\n", status.EnvTokenMasked)
	} else {
		fmt.Fprintln(w, "GITHUB_TOKEN:  not set")
	}

	fmt.Fprintf(w, "API base URL:  %s\n", status.APIBaseURL)

	if status.APIAuthenticated {
		fmt.Fprintf(w, "API access:    authenticated (via %s) as %s\n", status.APIAuthMethod, status.APIUsername)
	} else {
		fmt.Fprintln(w, "API access:    unauthenticated (public repos only)")
	}

	fmt.Fprintln(w)
	fmt.Fprintln(w, "Rate Limit")
	fmt.Fprintln(w, "──────────")
	fmt.Fprintf(w, "Limit:         %d\n", status.RateLimit)
	fmt.Fprintf(w, "Remaining:     %d\n", status.RateRemaining)
	fmt.Fprintf(w, "Resets at:     %s\n", status.RateReset.Format(time.RFC3339))
}