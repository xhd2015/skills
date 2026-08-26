## Expected

- Non-zero exit code.
- stderr contains `unknown command`.

## Side Effects

- None.

## Errors

- Process exits with error status; `Run` returns nil (exit error swallowed).

## Exit Code

- 1

```go
import (
	"github.com/xhd2015/doctest/session"
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if resp.ExitCode == 0 {
		t.Fatalf("expected non-zero exit, got 0\nstdout:\n%s\nstderr:\n%s", resp.Stdout, resp.Stderr)
	}
	combined := resp.Stdout + resp.Stderr
	if !strings.Contains(combined, "unknown command") {
		t.Fatalf("error missing unknown command:\n%s", combined)
	}
}
```