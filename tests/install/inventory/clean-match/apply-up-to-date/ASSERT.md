## Expected

- No error.
- stdout is exactly `Skill is up to date: <absDir>\n` (trailing newline).
- stdout does not contain `create:`, `update:`, `delete:`, or `Update skill at`.

## Side Effects

- `SKILL.md` content unchanged; no new files under the skill dir.

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
	assert.Output(t, resp.Stdout, fmt.Sprintf(`---
version: 2
---
Skill is up to date: %s
`, skillDir))
	for _, needle := range []string{"create:", "update:", "delete:", "Update skill at", "Installed skill to:"} {
		if strings.Contains(resp.Stdout, needle) {
			t.Fatalf("up-to-date stdout must not contain %q:\n%s", needle, resp.Stdout)
		}
	}
	got, readErr := os.ReadFile(filepath.Join(skillDir, "SKILL.md"))
	if readErr != nil {
		t.Fatalf("read SKILL.md: %v", readErr)
	}
	if string(got) != "# test skill\n" {
		t.Fatalf("SKILL.md content changed: %q", string(got))
	}
}
```
