## Expected

- Command exits with code 0.
- stdout contains `Installed skill to:`.
- `SKILL.md` exists at `.agents/skills/go-best-practice/SKILL.md`.
- `topics/skill-cli.md` exists under the same skill directory.

## Side Effects

- Skill directory and topic files are created on disk.

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
	topicMD := filepath.Join(skillDir, "topics", "skill-cli.md")

	if _, statErr := os.Stat(skillMD); statErr != nil {
		t.Fatalf("SKILL.md missing at %s: %v", skillMD, statErr)
	}
	if _, statErr := os.Stat(topicMD); statErr != nil {
		t.Fatalf("topics/skill-cli.md missing at %s: %v", topicMD, statErr)
	}
}
```