## Expected
- The command fails (non-zero exit code).
- The error message contains "no workflow runs matching" and "test".
- The error message lists available workflows: "lint" and "build".

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

	if !strings.Contains(combined, "no workflow runs matching") {
		t.Fatalf("error missing 'no workflow runs matching':\n%s", combined)
	}
	if !strings.Contains(combined, "test") {
		t.Fatalf("error should mention filter name 'test':\n%s", combined)
	}
	if !strings.Contains(combined, "lint") {
		t.Fatalf("error should list available workflow 'lint':\n%s", combined)
	}
	if !strings.Contains(combined, "build") {
		t.Fatalf("error should list available workflow 'build':\n%s", combined)
	}
}
```
