---
label: slow
explanation: nested lib receives explicit page param; launches Chromium
---

## Expected

- Exit code 0.
- stdout contains `explicit-page-ok`.

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
	if !strings.Contains(resp.Stdout, "explicit-page-ok") {
		t.Fatalf("stdout missing explicit-page-ok:\n%s", resp.Stdout)
	}
}
```