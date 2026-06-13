## Expected
- The command succeeds (exit code 0).
- stdout contains the "test" run with its failure status.
- stdout does NOT contain "lint" (filtered out).
- stdout contains the detected GitHub URL.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected exit code 0, got %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}

	stdout := resp.Stdout

	if !strings.Contains(stdout, "test") {
		t.Fatalf("stdout missing filtered run 'test':\n%s", stdout)
	}
	if strings.Contains(stdout, "lint") {
		t.Fatalf("stdout should NOT contain unfiltered run 'lint':\n%s", stdout)
	}
	if !strings.Contains(stdout, "failure") {
		t.Fatalf("stdout missing 'failure' conclusion:\n%s", stdout)
	}
	if !strings.Contains(stdout, "https://github.com/testowner/testrepo") {
		t.Fatalf("stdout missing detected GitHub URL:\n%s", stdout)
	}
}
```
