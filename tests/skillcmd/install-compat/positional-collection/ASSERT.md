## Expected

- Installs to `vendor/skills/demo-skill/SKILL.md`.
- Does not write collection-root `SKILL.md`.

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
	p := filepath.Join(resp.WorkDir, "vendor", "skills", "demo-skill", "SKILL.md")
	data, readErr := os.ReadFile(p)
	if readErr != nil {
		t.Fatalf("SKILL.md missing: %v\nstdout=%s", readErr, resp.Stdout)
	}
	if !strings.Contains(string(data), "installed via positional") {
		t.Fatalf("unexpected content: %q", data)
	}
	if _, e := os.Stat(filepath.Join(resp.WorkDir, "vendor", "skills", "SKILL.md")); e == nil {
		t.Fatal("must not write SKILL.md into the skills collection root")
	}
}
```
