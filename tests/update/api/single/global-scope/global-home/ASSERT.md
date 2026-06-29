## Expected

- stdout contains `Update skill at`.
- Global `SKILL.md` under `$HOME/.agents/skills/skill-alpha` is canonical.
- Project-local `.agents/skills/skill-alpha` does not exist.

## Side Effects

- Only the global install location is updated.

## Errors

- None.

```go
import (
	"os"
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
	if !strings.Contains(resp.Stdout, "Update skill at") {
		t.Fatalf("stdout missing update:\n%s", resp.Stdout)
	}
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("home dir: %v", err)
	}
	globalSkill := filepath.Join(home, skillAgentsDir("skill-alpha"), "SKILL.md")
	if got := readFile(t, globalSkill); got != skillAlphaContent {
		t.Fatalf("global SKILL.md not restored:\n%s", got)
	}
	localAgents := skillAgentsDir("skill-alpha")
	if pathExists(t, localAgents) {
		t.Fatalf("local %s should not exist", localAgents)
	}
}
```