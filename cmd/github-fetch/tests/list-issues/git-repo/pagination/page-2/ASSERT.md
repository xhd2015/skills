## Expected
- Exit code 0 with `Page 2 (2 items)`, issues #2 and #1, not #4/#3.

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
	stdout := resp.Stdout
	for _, want := range []string{"Page 2 (2 items)", "#2", "#1"} {
		if !strings.Contains(stdout, want) {
			t.Fatalf("stdout missing %q:\n%s", want, stdout)
		}
	}
	for _, absent := range []string{"#4", "#3"} {
		if strings.Contains(stdout, absent) {
			t.Fatalf("stdout should not contain %q:\n%s", absent, stdout)
		}
	}
}
```