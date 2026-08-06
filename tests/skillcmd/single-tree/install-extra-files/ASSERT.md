## Expected

- No error.
- stdout contains install confirmation (`Installed skill to:` for a fresh dir).
- On disk under `.agents/skills/demo-skill/`:
  - root `SKILL.md` exists
  - nested `skill-cli/TOPIC.md` exists
  - nested `a/b/TOPIC.md` exists (also in tree)
  - nested `skill-cli/SKILL.md` does **not** exist
  - legacy `topics/skill-cli.md` does **not** exist

## Side Effects

- Skill directory created with root + nested TOPIC.md extras only.

## Errors

- None.

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
	skillDir := filepath.Join(resp.WorkDir, ".agents", "skills", "demo-skill")
	skillMD := filepath.Join(skillDir, "SKILL.md")
	nestedTopic := filepath.Join(skillDir, "skill-cli", "TOPIC.md")
	nestedAB := filepath.Join(skillDir, "a", "b", "TOPIC.md")
	nestedOldSkill := filepath.Join(skillDir, "skill-cli", "SKILL.md")
	legacy := filepath.Join(skillDir, "topics", "skill-cli.md")
	if _, err := os.Stat(skillMD); err != nil {
		t.Fatalf("SKILL.md missing: %v", err)
	}
	if _, err := os.Stat(nestedTopic); err != nil {
		t.Fatalf("skill-cli/TOPIC.md missing: %v", err)
	}
	if _, err := os.Stat(nestedAB); err != nil {
		t.Fatalf("a/b/TOPIC.md missing: %v", err)
	}
	if _, err := os.Stat(nestedOldSkill); err == nil {
		t.Fatalf("skill-cli/SKILL.md must not be created")
	}
	if _, err := os.Stat(legacy); err == nil {
		t.Fatalf("legacy topics/skill-cli.md must not be created")
	}
}
```
