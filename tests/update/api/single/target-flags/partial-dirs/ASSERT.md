## Expected

- stdout mentions the codex absolute path (via `Update skill at` or `Skill is up to date` after restore).
- stdout has exactly one skill status line for this single skill update (one target processed).
- `.opencode/skills/skill-alpha` does not exist.

## Side Effects

- Codex `SKILL.md` restored to canonical content.

## Errors

- None.

```go
import (
	"path/filepath"
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if resp.Error != "" {
		t.Fatalf("expected no error, got: %s", resp.Error)
	}
	codexDir := skillCodexDir("skill-alpha")
	if !strings.Contains(resp.Stdout, "skill-alpha") && !strings.Contains(resp.Stdout, ".codex") {
		t.Fatalf("stdout should reference codex install, got:\n%s", resp.Stdout)
	}
	lines := strings.Split(strings.TrimSpace(resp.Stdout), "\n")
	statusLines := 0
	for _, line := range lines {
		if strings.Contains(line, "Skill is up to date") || strings.Contains(line, "Update skill at") || strings.Contains(line, "[dry-run]") {
			statusLines++
		}
	}
	if statusLines != 1 {
		t.Fatalf("expected 1 status line for single installed target, got %d:\n%s", statusLines, resp.Stdout)
	}
	opencodeDir := filepath.Join(resp.WorkDir, skillOpencodeDir("skill-alpha"))
	if pathExists(t, opencodeDir) {
		t.Fatalf("expected %s not to exist", opencodeDir)
	}
	codexSkill := filepath.Join(resp.WorkDir, skillMDPath(codexDir))
	if got := readFile(t, codexSkill); got != skillAlphaContent {
		t.Fatalf("codex SKILL.md not restored:\n%s", got)
	}
}
```