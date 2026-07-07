---
label: slow
explanation: eval mode launches playwright and forwards args to node -e subprocess
---

## Expected

- Exit code 0.
- stdout trimmed equals `["baz"]`.

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
	want := `["baz"]`
	got := strings.TrimSpace(resp.Stdout)
	if got != want {
		t.Fatalf("stdout = %q, want %q", got, want)
	}
}
```