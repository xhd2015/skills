## Expected
- Exit code 0 with message like `no open issues`.

## Exit Code
- 0

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected exit code 0, got %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	lower := strings.ToLower(resp.Stdout)
	if !strings.Contains(lower, "no") || !strings.Contains(lower, "issue") {
		t.Fatalf("stdout should indicate no issues:\n%s", resp.Stdout)
	}
}
```