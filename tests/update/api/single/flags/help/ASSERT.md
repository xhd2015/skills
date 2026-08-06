## Expected

- stdout documents `--dry-run` and `--global` (update flag surface).
- Usage mentions update as a subcommand-style command (`skills update` / `Usage:`),
  not a primary `--update` install flag.
- Help does not create skill directories.

## Side Effects

- No skill directories created.

## Errors

- None.

```go
import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if resp.Error != "" {
		t.Fatalf("expected no error, got: %s", resp.Error)
	}
	for _, needle := range []string{"--dry-run", "--global", "--color", "--no-color", "Usage:"} {
		if !strings.Contains(resp.Stdout, needle) {
			t.Fatalf("stdout missing %q:\n%s", needle, resp.Stdout)
		}
	}
	// Subcommand product surface: usage string carries "update", not --update as primary.
	if !strings.Contains(resp.Stdout, "update") {
		t.Fatalf("stdout should mention update usage:\n%s", resp.Stdout)
	}
	if strings.Contains(resp.Stdout, "--update") {
		t.Fatalf("help must not advertise --update flag as primary surface:\n%s", resp.Stdout)
	}
	if pathExists(t, filepath.Join(resp.WorkDir, skillAgentsDir("skill-alpha"))) {
		t.Fatalf("help must not create install dirs")
	}
}
```
