## Expected

- Command exits with code 0.
- stdout contains `Installed skill to:`.
- `SKILL.md` exists at `.agents/skills/go-best-practice/SKILL.md`.
- `skill-cli/SKILL.md` exists under the same skill directory.
- Legacy `topics/skill-cli.md` must not be required (nested layout replaces topics).

## Side Effects

- Skill directory and nested skill files are created on disk.

## Errors

- No error from `Run`.

## Exit Code

- 0

```go
import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run failed: %v\nstderr:\n%s", err, resp.Stderr)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit code = %d, want 0\nstderr:\n%s", resp.ExitCode, resp.Stderr)
	}
	if !strings.Contains(resp.Stdout, "Installed skill to:") {
		t.Fatalf("stdout missing install confirmation:\n%s", resp.Stdout)
	}

	skillDir := filepath.Join(resp.WorkDir, ".agents", "skills", "go-best-practice")
	skillMD := filepath.Join(skillDir, "SKILL.md")
	nested := filepath.Join(skillDir, "skill-cli", "SKILL.md")
	legacyTopic := filepath.Join(skillDir, "topics", "skill-cli.md")

	if _, statErr := os.Stat(skillMD); statErr != nil {
		t.Fatalf("SKILL.md missing at %s: %v", skillMD, statErr)
	}
	if _, statErr := os.Stat(nested); statErr != nil {
		t.Fatalf("skill-cli/SKILL.md missing at %s: %v", nested, statErr)
	}
	if _, statErr := os.Stat(legacyTopic); statErr == nil {
		t.Fatalf("legacy topics/skill-cli.md must not be installed at %s", legacyTopic)
	}
}
```
