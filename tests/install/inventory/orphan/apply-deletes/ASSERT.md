## Expected

- No error.
- stdout is exactly:
  - `Update skill at <absDir>`
  - `  delete: <absDir>/orphan.txt`
  (with trailing newline; no create/update lines).

## Side Effects

- `orphan.txt` is gone; `SKILL.md` still matches plan.

## Errors

- None.

```go
import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/assert"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if resp.Error != "" {
		t.Fatalf("expected no error, got: %s", resp.Error)
	}
	skillDir := absSkillDir(t, resp.WorkDir, "example-skill")
	orphan := filepath.Join(skillDir, "orphan.txt")
	assert.Output(t, resp.Stdout, fmt.Sprintf(`---
version: 2
---
Update skill at %s
  delete: %s
`, skillDir, orphan))
	if strings.Contains(resp.Stdout, "Skill is up to date") {
		t.Fatalf("orphan must not be treated as up to date:\n%s", resp.Stdout)
	}
	if _, statErr := os.Stat(orphan); !os.IsNotExist(statErr) {
		t.Fatalf("orphan.txt should be deleted; stat err=%v", statErr)
	}
	got, readErr := os.ReadFile(filepath.Join(skillDir, "SKILL.md"))
	if readErr != nil {
		t.Fatalf("read SKILL.md: %v", readErr)
	}
	if string(got) != "# test skill\n" {
		t.Fatalf("SKILL.md content mismatch: %q", string(got))
	}
}
```
