---
label: slow
explanation: bare file alias routes to file mode and forwards script args
---

## Expected

- Exit code 0.
- stdout trimmed equals `["foo","bar"]`.

## Exit Code

- 0

```go
import (
	"strings"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run failed: %v\nstderr:\n%s", err, resp.Stderr)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit code = %d, want 0\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	want := `["foo","bar"]`
	got := strings.TrimSpace(resp.Stdout)
	if got != want {
		t.Fatalf("stdout = %q, want %q", got, want)
	}
}
```