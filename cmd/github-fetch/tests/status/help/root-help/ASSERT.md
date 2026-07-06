## Expected

- Exit code 0.
- stdout root help lists the `status` command alongside other sub-commands.

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

	if !strings.Contains(resp.Stdout, "status") {
		t.Fatalf("root help should list status command, got:\n%s", resp.Stdout)
	}
	if !strings.Contains(strings.ToLower(resp.Stdout), "commands") {
		t.Fatalf("root help should include commands section, got:\n%s", resp.Stdout)
	}
}
```