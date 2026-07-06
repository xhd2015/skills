## Expected Output

```text
---
version: 2
---
Auth Status
───────────
gh CLI:        available
GitHub host:   github.com
Logged in as:  ghuser
Token scopes:  repo
GITHUB_TOKEN:  set (gho_****)
API base URL:  __API_BASE_URL__
API access:    authenticated (via GITHUB_TOKEN) as testuser

Rate Limit
──────────
Limit:         5000
Remaining:     4987
Resets at:     .+
```

## Expected

- Exit code 0.
- stdout shows GITHUB_TOKEN set with masked value and authenticated via GITHUB_TOKEN (not gh).

## Side Effects

- API probes use Bearer env token, not `gh auth token` output.

## Errors

- No error from `Run`.

## Exit Code

- 0

```go
import (
	"testing"

	"github.com/xhd2015/doctest/assert"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected exit code 0, got %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}

	assert.Output(t, resp.Stdout, `---
version: 2
__API_BASE_URL__: type=string, example=http://127.0.0.1:54321, mock API base URL
---
Auth Status
───────────
gh CLI:        available
GitHub host:   github.com
Logged in as:  ghuser
Token scopes:  repo
GITHUB_TOKEN:  set (gho_****)
API base URL:  __API_BASE_URL__
API access:    authenticated (via GITHUB_TOKEN) as testuser

Rate Limit
──────────
Limit:         5000
Remaining:     4987
Resets at:     .+
`)
}
```