## Expected

- No error.
- stdout:
  - `Installed skill to: <absDir>`
  - `  create: <absDir>/SKILL.md`
  - `  create: <absDir>/nested/TOPIC.md`
  (creates sorted by relative path: `SKILL.md` before `nested/TOPIC.md`).
- Does not use `Update skill at` header.

## Side Effects

- Both planned files exist with planned content.
- No unplanned files required.

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
	skillMD := filepath.Join(skillDir, "SKILL.md")
	nested := filepath.Join(skillDir, "nested", "TOPIC.md")
	assert.Output(t, resp.Stdout, fmt.Sprintf(`---
version: 3
---
Installed skill to: %s
  create: %s
  create: %s
`, skillDir, skillMD, nested))
	if strings.Contains(resp.Stdout, "Update skill at") {
		t.Fatalf("fresh install must not use Update header:\n%s", resp.Stdout)
	}
	if strings.Contains(resp.Stdout, "update:") || strings.Contains(resp.Stdout, "delete:") {
		t.Fatalf("fresh install should only create:\n%s", resp.Stdout)
	}
	gotRoot, err1 := os.ReadFile(skillMD)
	gotNested, err2 := os.ReadFile(nested)
	if err1 != nil || err2 != nil {
		t.Fatalf("planned files missing: root=%v nested=%v", err1, err2)
	}
	if string(gotRoot) != "# test skill\n" {
		t.Fatalf("SKILL.md content = %q", string(gotRoot))
	}
	if string(gotNested) != "# nested topic\n" {
		t.Fatalf("nested/TOPIC.md content = %q", string(gotNested))
	}
}
```
