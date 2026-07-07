---
label: slow
explanation: launches Chromium; script --help must forward past CLI boundary
---

## Expected

- Exit code 0.
- stdout contains script marker `SCRIPT_HELP_OK`.
- stdout does **not** contain CLI help (`Commands:`).

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
	if !strings.Contains(out, "SCRIPT_HELP_OK") {
		t.Fatalf("stdout missing SCRIPT_HELP_OK:\n%s", out)
	}
	if strings.Contains(out, "Commands:") {
		t.Fatalf("stdout must not contain CLI help (Commands:), got:\n%s", out)
	}
}
```