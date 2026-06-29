## Expected

- stdout contains `Skill is up to date`.
- stdout does not contain `Update skill at`.

## Side Effects

- `SKILL.md` content remains canonical.

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
	if !strings.Contains(resp.Stdout, "Skill is up to date") {
		t.Fatalf("stdout missing up-to-date message:\n%s", resp.Stdout)
	}
	if strings.Contains(resp.Stdout, "Update skill at") {
		t.Fatalf("stdout should not report update:\n%s", resp.Stdout)
	}
	skillPath := filepath.Join(resp.WorkDir, skillMDPath(skillAgentsDir("skill-alpha")))
	if got := readFile(t, skillPath); got != skillAlphaContent {
		t.Fatalf("SKILL.md content mismatch:\n%s", got)
	}
}
```