## Expected

- Non-zero exit code.
- stderr contains an error message about API probe failure.

## Side Effects

- Mock server returns 500 for API requests.

## Errors

- Command exits with failure; `Run` itself succeeds in capturing output.

## Exit Code

- Non-zero

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
	if resp.ExitCode == 0 {
		t.Fatalf("expected non-zero exit code, got 0\nstdout:\n%s\nstderr:\n%s", resp.Stdout, resp.Stderr)
	}

	stderr := strings.ToLower(resp.Stderr)
	if stderr == "" {
		t.Fatalf("expected error on stderr, got empty stderr\nstdout:\n%s", resp.Stdout)
	}
	if strings.Contains(stderr, "unknown command") {
		t.Fatalf("expected API probe failure, not missing sub-command:\n%s", resp.Stderr)
	}
	if !strings.Contains(stderr, "error") && !strings.Contains(stderr, "500") && !strings.Contains(stderr, "api") && !strings.Contains(stderr, "rate") {
		t.Fatalf("stderr should describe API failure, got:\n%s", resp.Stderr)
	}
}
```