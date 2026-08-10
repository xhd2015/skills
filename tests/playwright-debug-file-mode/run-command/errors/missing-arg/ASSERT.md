## Expected

- Non-zero exit code.
- Error mentions a `.js` file is required.

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
	if !strings.Contains(combined, "file") && !strings.Contains(combined, ".js") {
		t.Fatalf("error should mention file/.js, got:\n%s", combined)
	}
	if !strings.Contains(combined, "require") {
		t.Fatalf("error should mention required file, got:\n%s", combined)
	}
}
```