## Expected

- Exit code 0.
- stdout contains usage/help text for the `status` sub-command.

## Side Effects

- None beyond printing help.

## Errors

- No error from `Run`.

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

	stdout := strings.ToLower(resp.Stdout)
	if !strings.Contains(stdout, "status") {
		t.Fatalf("status help should mention status, got:\n%s", resp.Stdout)
	}
	if !strings.Contains(stdout, "usage") && !strings.Contains(stdout, "help") {
		t.Fatalf("status help should look like usage text, got:\n%s", resp.Stdout)
	}
}
```