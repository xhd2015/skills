## Expected

- No error.
- stdout:
  - `Update skill at <absDir>`
  - `  update: <absDir>/SKILL.md`
- No create/delete lines.

## Side Effects

- SKILL.md content is the new plan content.
- Skill dir is not wholesale removed (other unrelated structure would stay; none here).

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
	skillMD := filepath.Join(skillDir, "SKILL.md")
	assert.Output(t, resp.Stdout, fmt.Sprintf(`---
version: 3
---
Update skill at %s
  update: %s
`, skillDir, skillMD))
	if strings.Contains(resp.Stdout, "create:") || strings.Contains(resp.Stdout, "delete:") {
		t.Fatalf("content-only update should not create/delete:\n%s", resp.Stdout)
	}
	got, readErr := os.ReadFile(skillMD)
	if readErr != nil {
		t.Fatalf("read SKILL.md: %v", readErr)
	}
	if string(got) != "# new skill\n" {
		t.Fatalf("SKILL.md content = %q, want new skill", string(got))
	}
}
```
