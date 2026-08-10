## Expected

- Non-zero exit code.
- Error mentions file not found or requires existing `.js` file.

## Exit Code

- 1

```go
import (
	"strings"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if resp.ExitCode == 0 {
		t.Fatalf("expected non-zero exit, got 0\nstdout:\n%s\nstderr:\n%s", resp.Stdout, resp.Stderr)
	}
	combined := strings.ToLower(resp.Stdout + resp.Stderr)
	if !strings.Contains(combined, "missing.js") && !strings.Contains(combined, "not found") && !strings.Contains(combined, "no such") {
		t.Fatalf("error should mention missing file, got:\n%s", combined)
	}
}
```