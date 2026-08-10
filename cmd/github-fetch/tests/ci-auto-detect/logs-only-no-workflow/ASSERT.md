## Expected
- The command succeeds (exit code 0).
- stdout contains a status header for the latest run: `Workflow: lint (Run #200`.
- stdout contains the log content: `linter output`.
- stdout contains the detected GitHub URL.

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

	if !strings.Contains(stdout, "Workflow: lint") {
		t.Fatalf("stdout missing workflow name 'lint' for latest run:\n%s", stdout)
	}
	if !strings.Contains(stdout, "Run #200") {
		t.Fatalf("stdout missing run ID #200:\n%s", stdout)
	}
	if !strings.Contains(stdout, "failure") {
		t.Fatalf("stdout missing conclusion 'failure':\n%s", stdout)
	}
	if !strings.Contains(stdout, "linter output") {
		t.Fatalf("stdout missing log content:\n%s", stdout)
	}
	if !strings.Contains(stdout, "https://github.com/testowner/testrepo") {
		t.Fatalf("stdout missing detected GitHub URL:\n%s", stdout)
	}
}
```
