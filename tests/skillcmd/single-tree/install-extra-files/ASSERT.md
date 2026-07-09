## Expected

- Scenario-specific success or error as asserted.

## Side Effects

- install-extra-files creates nested skill-cli/SKILL.md under the skill dir.

## Errors

- reject-dotdot expects a non-empty error.

```go
import (
	"path/filepath"
	"os"
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
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
	nested := filepath.Join(skillDir, "skill-cli", "SKILL.md")
	legacy := filepath.Join(skillDir, "topics", "skill-cli.md")
	if _, err := os.Stat(skillMD); err != nil {
		t.Fatalf("SKILL.md missing: %v", err)
	}
	if _, err := os.Stat(nested); err != nil {
		t.Fatalf("skill-cli/SKILL.md missing: %v", err)
	}
	if _, err := os.Stat(legacy); err == nil {
		t.Fatalf("legacy topics/skill-cli.md must not be created")
	}
}
```
