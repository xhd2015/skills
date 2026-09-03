## Expected

- No error.
- Skill root is `out/demo-skill/SKILL.md` (no extra nesting).

```go
import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if resp.Error != "" {
		t.Fatalf("unexpected error: %s", resp.Error)
	}
	p := filepath.Join(resp.WorkDir, "out", "demo-skill", "SKILL.md")
	data, readErr := os.ReadFile(p)
	if readErr != nil {
		t.Fatalf("SKILL.md missing: %v\nstdout=%s", readErr, resp.Stdout)
	}
	if !strings.Contains(string(data), "installed via --dir") {
		t.Fatalf("unexpected content: %q", data)
	}
	nested := filepath.Join(resp.WorkDir, "out", "demo-skill", "demo-skill", "SKILL.md")
	if _, e := os.Stat(nested); e == nil {
		t.Fatal("must not double-nest when basename matches skill name")
	}
}
```
