## Expected

- No error.
- stdout order: header, then `create:` (TOPIC), then `delete:` (old SKILL):
  - `Update skill at <absDir>`
  - `  create: <absDir>/a/TOPIC.md`
  - `  delete: <absDir>/a/SKILL.md`
- No separate `rename:` op; no up-to-date claim.

## Side Effects

- `a/TOPIC.md` exists with nested body.
- `a/SKILL.md` is gone.
- Root SKILL.md unchanged.

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
	oldNested := filepath.Join(skillDir, "a", "SKILL.md")
	newNested := filepath.Join(skillDir, "a", "TOPIC.md")
	assert.Output(t, resp.Stdout, fmt.Sprintf(`---
version: 2
---
Update skill at %s
  create: %s
  delete: %s
`, skillDir, newNested, oldNested))
	if strings.Contains(resp.Stdout, "Skill is up to date") {
		t.Fatalf("rename must not be up to date:\n%s", resp.Stdout)
	}
	if strings.Contains(resp.Stdout, "rename:") {
		t.Fatalf("must not emit rename: op:\n%s", resp.Stdout)
	}
	if _, statErr := os.Stat(oldNested); !os.IsNotExist(statErr) {
		t.Fatalf("a/SKILL.md should be deleted; stat err=%v", statErr)
	}
	got, readErr := os.ReadFile(newNested)
	if readErr != nil {
		t.Fatalf("a/TOPIC.md missing: %v", readErr)
	}
	if string(got) != "# nested topic body\n" {
		t.Fatalf("a/TOPIC.md content = %q", string(got))
	}
	root, readErr := os.ReadFile(filepath.Join(skillDir, "SKILL.md"))
	if readErr != nil {
		t.Fatalf("read root SKILL.md: %v", readErr)
	}
	if string(root) != "# test skill\n" {
		t.Fatalf("root SKILL.md mismatch: %q", string(root))
	}
}
```
