## Expected

- Non-zero exit code.
- stderr contains `unknown skill sub-command`.
- stderr mentions `skill show` or `skill install`.

## Side Effects

- None.

## Errors

- Process exits with error status; `Run` returns nil (exit error swallowed).

## Exit Code

- 1

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if resp.ExitCode == 0 {
		t.Fatalf("expected non-zero exit, got 0\nstdout:\n%s\nstderr:\n%s", resp.Stdout, resp.Stderr)
	}
	combined := resp.Stdout + resp.Stderr
	if !strings.Contains(combined, "unknown skill sub-command") {
		t.Fatalf("error missing unknown skill sub-command:\n%s", combined)
	}
	if !strings.Contains(combined, "skill show") && !strings.Contains(combined, "skill install") {
		t.Fatalf("error must mention skill show or skill install:\n%s", combined)
	}
}
```