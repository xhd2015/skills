## Expected

- No error.
- stdout contains `Installed skill to:`.
- File `.agents/skills/demo-skill/SKILL.md` exists with installed content.

## Side Effects

- Skill directory created under workdir.

```go
import (
	"github.com/xhd2015/doctest/session"
	"os"
	"path/filepath"
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
	if !strings.Contains(resp.Stdout, "Installed skill to:") {
		t.Fatalf("stdout missing install confirmation:\n%s", resp.Stdout)
	}
	p := filepath.Join(resp.WorkDir, ".agents", "skills", "demo-skill", "SKILL.md")
	data, readErr := os.ReadFile(p)
	if readErr != nil {
		t.Fatalf("SKILL.md missing: %v", readErr)
	}
	if !strings.Contains(string(data), "installed via skillcmd") {
		t.Fatalf("unexpected content: %q", data)
	}
}
```
