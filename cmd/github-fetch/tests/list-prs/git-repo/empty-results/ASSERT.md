## Expected
- Exit code 0.
- stdout contains a message like `no open pull requests` (case-insensitive match on key words).

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
		t.Fatalf("run failed: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected exit code 0, got %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	lower := strings.ToLower(resp.Stdout)
	if !strings.Contains(lower, "no") || !strings.Contains(lower, "pull") {
		t.Fatalf("stdout should indicate no pull requests:\n%s", resp.Stdout)
	}
}
```