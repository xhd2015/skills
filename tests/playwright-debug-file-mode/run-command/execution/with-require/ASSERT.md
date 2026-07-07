---
label: slow
explanation: file mode with createRequire relative to script dir
---

## Expected

- Exit code 0.
- stdout contains `require-ok`.

## Exit Code

- 0

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run failed: %v\nstderr:\n%s", err, resp.Stderr)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit code = %d, want 0\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	if !strings.Contains(resp.Stdout, "require-ok") {
		t.Fatalf("stdout missing require-ok:\n%s", resp.Stdout)
	}
}
```