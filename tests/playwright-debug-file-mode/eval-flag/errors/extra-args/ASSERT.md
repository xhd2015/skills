## Expected

- Non-zero exit code.
- Error mentions unexpected arguments after `--eval`.

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
	combined := strings.ToLower(resp.Stdout + resp.Stderr)
	if !strings.Contains(combined, "unexpected") && !strings.Contains(combined, "extra") && !strings.Contains(combined, "argument") {
		t.Fatalf("error should mention unexpected/extra arguments, got:\n%s", combined)
	}
}
```