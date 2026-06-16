## Expected
- The command fails (non-zero exit code).
- The error message indicates that no origin remote is configured.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}
	if resp.ExitCode == 0 {
		t.Fatalf("expected non-zero exit code, got 0\nstdout:\n%s\nstderr:\n%s", resp.Stdout, resp.Stderr)
	}
	combined := resp.Stdout + resp.Stderr
	if !strings.Contains(strings.ToLower(combined), "origin") &&
		!strings.Contains(strings.ToLower(combined), "remote") {
		t.Fatalf("error should mention origin/remote, got:\n%s", combined)
	}
}
```
