## Expected
- The command succeeds (exit code 0).
- stdout contains a status header: `Workflow: test (Run #12345, event: push, branch: main) — failed`.
- stdout contains the log content: `build step failed`.
- stdout contains the detected GitHub URL: `https://github.com/testowner/testrepo`.

## Side Effects
- The mock API received requests for repo info, workflow runs, workflow jobs, and job logs.

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

	stdout := resp.Stdout

	if !strings.Contains(stdout, "Workflow: test") {
		t.Fatalf("stdout missing workflow name 'test':\n%s", stdout)
	}
	if !strings.Contains(stdout, "Run #12345") {
		t.Fatalf("stdout missing run ID #12345:\n%s", stdout)
	}
	if !strings.Contains(stdout, "failure") {
		t.Fatalf("stdout missing conclusion 'failure':\n%s", stdout)
	}
	if !strings.Contains(stdout, "build step failed") {
		t.Fatalf("stdout missing log content:\n%s", stdout)
	}
	if !strings.Contains(stdout, "https://github.com/testowner/testrepo") {
		t.Fatalf("stdout missing detected GitHub URL:\n%s", stdout)
	}
}
```
