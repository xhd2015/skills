## Expected

- Non-zero exit code.
- Error mentions script argument is required.

## Exit Code

- 1

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if resp.ExitCode == 0 {
		t.Fatalf("expected non-zero exit, got 0\nstdout:\n%s\nstderr:\n%s", resp.Stdout, resp.Stderr)
	}
	combined := strings.ToLower(resp.Stdout + resp.Stderr)
	if !strings.Contains(combined, "script") && !strings.Contains(combined, "eval") {
		t.Fatalf("error should mention script/eval, got:\n%s", combined)
	}
	if !strings.Contains(combined, "require") {
		t.Fatalf("error should mention required argument, got:\n%s", combined)
	}
}
```