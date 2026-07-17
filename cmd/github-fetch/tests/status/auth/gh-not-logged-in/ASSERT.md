## Expected Output

```text
---
version: 3
---
Auth Status
───────────
gh CLI:        available
GITHUB_TOKEN:  not set
API base URL:  __API_BASE_URL__
API access:    unauthenticated \(public repos only\)

Rate Limit
──────────
Limit:         5000
Remaining:     4987
Resets at:     .+
```

## Expected

- Exit code 0.
- stdout shows gh available but API access unauthenticated.

## Side Effects

- Fake gh `auth status` invoked; no gh token used for API.

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
version: 3
__API_BASE_URL__: type=string, example=http://127.0.0.1:54321, mock API base URL
---
Auth Status
───────────
gh CLI:        available
GITHUB_TOKEN:  not set
API base URL:  __API_BASE_URL__
API access:    unauthenticated \(public repos only\)

Rate Limit
──────────
Limit:         5000
Remaining:     4987
Resets at:     .+
`)
}
```