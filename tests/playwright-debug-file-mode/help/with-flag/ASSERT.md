## Expected

- Exit code 0.
- stdout contains `Usage:`.
- stdout mentions `run`, eval flags, and file alias (same as no-args help).
- stdout mentions trailing script arguments are passed through to `process.argv`.

## Side Effects

- None.

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
		t.Fatalf("Run failed: %v\nstderr:\n%s", err, resp.Stderr)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit code = %d, want 0\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	out := resp.Stdout
	if !strings.Contains(out, "Usage:") {
		t.Fatalf("stdout missing Usage:\n%s", out)
	}
	lower := strings.ToLower(out)
	if !strings.Contains(lower, "run") {
		t.Fatalf("stdout must mention run:\n%s", out)
	}
	if !strings.Contains(out, "-e") && !strings.Contains(out, "--eval") {
		t.Fatalf("stdout must mention eval flags:\n%s", out)
	}
	if !strings.Contains(lower, ".js") {
		t.Fatalf("stdout must mention .js:\n%s", out)
	}
	if !strings.Contains(out, "process.argv") {
		t.Fatalf("stdout must mention process.argv for script arg pass-through:\n%s", out)
	}
}
```