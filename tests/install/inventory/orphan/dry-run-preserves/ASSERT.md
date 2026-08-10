## Expected

- No error.
- Dry-run stdout reports update header + delete for orphan (each line prefixed
  with `[dry-run] `).
- stdout does not claim up to date.

## Side Effects

- `orphan.txt` remains on disk with original content.
- `SKILL.md` unchanged.

## Errors

- None.

```go
import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/session"

	"github.com/xhd2015/doctest/assert"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if resp.Error != "" {
		t.Fatalf("expected no error, got: %s", resp.Error)
	}
	skillDir := absSkillDir(t, resp.WorkDir, "example-skill")
	orphan := filepath.Join(skillDir, "orphan.txt")
	assert.Output(t, resp.Stdout, fmt.Sprintf(`---
version: 3
---
\[dry-run\] Update skill at %s
\[dry-run\]   delete: %s
`, skillDir, orphan))
	if strings.Contains(resp.Stdout, "Skill is up to date") {
		t.Fatalf("orphan dry-run must not be up to date:\n%s", resp.Stdout)
	}
	got, readErr := os.ReadFile(orphan)
	if readErr != nil {
		t.Fatalf("orphan.txt must still exist: %v", readErr)
	}
	if string(got) != "leftover\n" {
		t.Fatalf("orphan content changed: %q", string(got))
	}
}
```
