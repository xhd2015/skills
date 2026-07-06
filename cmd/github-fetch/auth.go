package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type AuthMethod string

const (
	AuthMethodEnvToken AuthMethod = "GITHUB_TOKEN"
	AuthMethodGh       AuthMethod = "gh"
	AuthMethodNone     AuthMethod = "none"
)

var (
	ghLoggedInAsRe    = regexp.MustCompile(`(?i)logged in to .+ as (\S+)`)
	ghLoggedInAcctRe  = regexp.MustCompile(`(?i)logged in to (\S+) account (\S+)`)
	ghTokenScopesRe   = regexp.MustCompile(`(?i)^Token scopes:\s*(.+)$`)
	ghHostFromLoginRe = regexp.MustCompile(`(?i)logged in to (\S+)`)
)

func resolveAuthToken() (token string, method AuthMethod) {
	if token := strings.TrimSpace(os.Getenv("GITHUB_TOKEN")); token != "" {
		return token, AuthMethodEnvToken
	}
	if isGHAvailable() {
		out, err := runCmd("gh", "auth", "token")
		if err == nil {
			token = strings.TrimSpace(out)
			if token != "" {
				return token, AuthMethodGh
			}
		}
	}
	return "", AuthMethodNone
}

func authenticatedGet(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	return authenticatedDo(req)
}

func authenticatedDo(req *http.Request) (*http.Response, error) {
	if token, _ := resolveAuthToken(); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return http.DefaultClient.Do(req)
}

func maskToken(token string) string {
	if len(token) <= 4 {
		return "****"
	}
	return token[:4] + "****"
}

func getGhAuthInfo() (available, loggedIn bool, host, username string, scopes []string) {
	if !isGHAvailable() {
		return false, false, "", "", nil
	}
	available = true

	cmd := exec.Command("gh", "auth", "status")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return available, false, "", "", nil
	}

	host, username, scopes = parseGhAuthStatusOutput(string(out))
	return available, true, host, username, scopes
}

func parseGhAuthStatusOutput(output string) (host, username string, scopes []string) {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		line = strings.TrimPrefix(line, "✓ ")
		line = strings.TrimSpace(line)

		if m := ghLoggedInAcctRe.FindStringSubmatch(line); m != nil {
			host = m[1]
			username = m[2]
			continue
		}
		if m := ghLoggedInAsRe.FindStringSubmatch(line); m != nil {
			username = m[1]
			if host == "" {
				if hm := ghHostFromLoginRe.FindStringSubmatch(line); hm != nil {
					host = hm[1]
				}
			}
			continue
		}
		if m := ghTokenScopesRe.FindStringSubmatch(line); m != nil {
			scopes = parseGhTokenScopes(m[1])
			continue
		}
		if host == "" && !strings.Contains(strings.ToLower(line), "token scopes") &&
			!strings.Contains(strings.ToLower(line), "logged in") {
			if strings.Contains(line, ".") || line == "github.com" {
				host = line
			}
		}
	}
	return host, username, scopes
}

func parseGhTokenScopes(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	scopes := make([]string, 0, len(parts))
	for _, part := range parts {
		scope := strings.TrimSpace(part)
		scope = strings.Trim(scope, "'\"")
		if scope != "" {
			scopes = append(scopes, scope)
		}
	}
	return scopes
}