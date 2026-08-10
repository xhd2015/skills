## Expected
- The command succeeds (exit code 0).
- stdout shows logs for run #200 (the lint workflow).
- stdout does NOT show logs for run #100.
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
		t.Fatalf("stdout missing workflow name 'lint':\n%s", stdout)
	}
	if !strings.Contains(stdout, "Run #200") {
		t.Fatalf("stdout missing run ID #200:\n%s", stdout)
	}
	if !strings.Contains(stdout, "linter error: unused variable") {
		t.Fatalf("stdout missing log content:\n%s", stdout)
	}
	if strings.Contains(stdout, "Run #100") {
		t.Fatalf("stdout should NOT contain run #100:\n%s", stdout)
	}
	if !strings.Contains(stdout, "https://github.com/testowner/testrepo") {
		t.Fatalf("stdout missing detected GitHub URL:\n%s", stdout)
	}
}
```
