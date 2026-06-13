## Expected
- The command fails (non-zero exit code).
- The error message indicates that the user is not in a git repository.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}
	// The Run function should NOT return an error for non-zero exit;
	// it should capture the exit code in resp.ExitCode.
	if resp.ExitCode == 0 {
		t.Fatalf("expected non-zero exit code, got 0\nstdout:\n%s\nstderr:\n%s", resp.Stdout, resp.Stderr)
	}
	combined := resp.Stdout + resp.Stderr
	if !strings.Contains(strings.ToLower(combined), "git") &&
		!strings.Contains(strings.ToLower(combined), "repo") {
		t.Fatalf("error should mention git/repo, got:\n%s", combined)
	}
}
```
