## Expected
- The command fails with a non-zero exit code.
- stderr mentions that the user is not in a git repository (or cannot auto-detect repo).

## Exit Code
- Non-zero

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
	combined := strings.ToLower(resp.Stdout + resp.Stderr)
	if !strings.Contains(combined, "git") && !strings.Contains(combined, "repo") {
		t.Fatalf("error should mention git/repo, got:\n%s", combined)
	}
}
```