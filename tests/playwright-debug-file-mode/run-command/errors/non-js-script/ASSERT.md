## Expected

- Non-zero exit code.
- Error states `run` requires an **existing** `.js` file (not eval fallback).

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
	if !strings.Contains(combined, ".js") {
		t.Fatalf("error should mention .js file, got:\n%s", combined)
	}
	if !strings.Contains(combined, "exist") && !strings.Contains(combined, "existing") {
		t.Fatalf("error should mention existing file, got:\n%s", combined)
	}
}
```