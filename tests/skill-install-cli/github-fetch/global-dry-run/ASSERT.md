## Expected

- Command exits with code 0.
- stdout contains `[dry-run]`.
- stdout mentions `HOME/.agents/skills/github-fetch`.

## Side Effects

- No files are created under the temp HOME directory.

## Errors

- No error from `Run`.

## Exit Code

- 0

```go
import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run failed: %v\nstderr:\n%s", err, resp.Stderr)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit code = %d, want 0\nstderr:\n%s", resp.ExitCode, resp.Stderr)
	}
	expected := filepath.Join(resp.HomeDir, ".agents", "skills", "github-fetch")
	if !strings.Contains(resp.Stdout, expected) {
		t.Fatalf("stdout missing global target %q:\n%s", expected, resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "[dry-run]") {
		t.Fatalf("stdout missing [dry-run]:\n%s", resp.Stdout)
	}
	if _, statErr := os.Stat(expected); !os.IsNotExist(statErr) {
		t.Fatalf("dry-run must not create global target dir %s", expected)
	}
}
```