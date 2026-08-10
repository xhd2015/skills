## Expected

- Exit code 0.
- stdout `API base URL` line equals the mock server URL from `GITHUB_API_BASE_URL`.
- stdout still includes rate-limit section from mock `/rate_limit`.

## Side Effects

- Mock server handles API requests against the custom base URL.

## Errors

- No error from `Run`.

## Exit Code

- 0

```go
import (
	"strings"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected exit code 0, got %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}

	if !strings.Contains(resp.Stdout, "API base URL:") {
		t.Fatalf("stdout missing API base URL line:\n%s", resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, resp.APIBaseURL) {
		t.Fatalf("stdout should show mock API base URL %q:\n%s", resp.APIBaseURL, resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "Limit:") || !strings.Contains(resp.Stdout, "5000") {
		t.Fatalf("stdout missing rate limit from mock API:\n%s", resp.Stdout)
	}
}
```