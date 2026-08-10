## Expected

- No error.
- stdout:
  - `Update skill at <absDir>`
  - `  delete: <absDir>/a/TOPIC.md`
- Not up to date.

## Side Effects

- `a/TOPIC.md` removed.
- Empty directory `a/` under the skill root is removed best-effort.
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
	nested := filepath.Join(skillDir, "a", "TOPIC.md")
	assert.Output(t, resp.Stdout, fmt.Sprintf(`---
version: 3
---
Update skill at %s
  delete: %s
`, skillDir, nested))
	if strings.Contains(resp.Stdout, "Skill is up to date") {
		t.Fatalf("dropped nested file must not be up to date:\n%s", resp.Stdout)
	}
	if _, statErr := os.Stat(nested); !os.IsNotExist(statErr) {
		t.Fatalf("a/TOPIC.md should be deleted; stat err=%v", statErr)
	}
	// best-effort empty dir cleanup
	if _, statErr := os.Stat(filepath.Join(skillDir, "a")); !os.IsNotExist(statErr) {
		t.Fatalf("empty dir a/ should be removed best-effort; stat err=%v", statErr)
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
