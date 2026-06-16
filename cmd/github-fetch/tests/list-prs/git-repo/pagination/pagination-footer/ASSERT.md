## Expected
- Exit code 0.
- stdout contains `Page 1 (30 items)`.
- stdout contains `More results available` and `--page 2`.

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
	for _, want := range []string{
		"Page 1 (30 items)",
		"More results available",
		"--page 2",
	} {
		if !strings.Contains(stdout, want) {
			t.Fatalf("stdout missing %q:\n%s", want, stdout)
		}
	}
}
```