## Expected

- No error.
- stdout contains `[dry-run]` and `.agents/skills/foo`.

## Side Effects

- No directory created under workdir.

```go
import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if resp.Error != "" {
		t.Fatalf("unexpected error: %s", resp.Error)
	}
	if !strings.Contains(resp.Stdout, "[dry-run]") {
		t.Fatalf("stdout missing [dry-run]:\n%s", resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, ".agents/skills/foo") {
		t.Fatalf("stdout missing foo target:\n%s", resp.Stdout)
	}
	target := filepath.Join(resp.WorkDir, ".agents", "skills", "foo")
	if _, statErr := os.Stat(target); !os.IsNotExist(statErr) {
		t.Fatalf("dry-run must not create %s", target)
	}
}
```
