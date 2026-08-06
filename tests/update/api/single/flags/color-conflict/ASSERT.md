## Expected

- Update fails with a clear mutual-exclusion error for `--color` and `--no-color`.
- No skill directories created.

## Errors

- Error message contains both flag names and "cannot be specified together" (or equivalent).

## Exit Code

non-zero product error (surfaced in Response.Error)

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
	if resp.Error == "" {
		t.Fatalf("expected color flag conflict error, got empty Error; stdout:\n%s", resp.Stdout)
	}
	lower := strings.ToLower(resp.Error)
	if !strings.Contains(lower, "--color") || !strings.Contains(lower, "--no-color") {
		t.Fatalf("error should mention both flags: %s", resp.Error)
	}
	if !strings.Contains(lower, "together") && !strings.Contains(lower, "cannot") {
		t.Fatalf("error should state mutual exclusion: %s", resp.Error)
	}
	if pathExists(t, filepath.Join(resp.WorkDir, skillAgentsDir("skill-alpha"))) {
		t.Fatalf("conflict must not create install dirs")
	}
}
```
