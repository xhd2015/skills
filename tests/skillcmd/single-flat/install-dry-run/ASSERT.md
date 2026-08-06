## Expected

- No error.
- See scenario-specific stdout checks in Assert.

## Side Effects

- install-dry-run must not create directories.

## Errors

- None.

```go
import (
	"github.com/xhd2015/doctest/session"
	"path/filepath"
	"os"
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if resp.Error != "" {
		t.Fatalf("unexpected error: %s", resp.Error)
	}
	if !strings.Contains(resp.Stdout, "[dry-run]") {
		t.Fatalf("stdout missing [dry-run]:\n%s", resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, ".agents/skills/demo-skill") {
		t.Fatalf("stdout missing default target:\n%s", resp.Stdout)
	}
	target := filepath.Join(resp.WorkDir, ".agents", "skills", "demo-skill")
	if _, statErr := os.Stat(target); !os.IsNotExist(statErr) {
		t.Fatalf("dry-run must not create %s", target)
	}
}
```
