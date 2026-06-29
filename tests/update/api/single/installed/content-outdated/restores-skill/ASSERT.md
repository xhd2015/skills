## Expected

- stdout contains `Update skill at`.
- `SKILL.md` on disk equals canonical embedded content.

## Side Effects

- Drifted content is replaced.

## Errors

- None.

```go
import (
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
		t.Fatalf("stdout missing update message:\n%s", resp.Stdout)
	}
	skillPath := filepath.Join(resp.WorkDir, skillMDPath(skillAgentsDir("skill-alpha")))
	if got := readFile(t, skillPath); got != skillAlphaContent {
		t.Fatalf("SKILL.md not restored, got:\n%s", got)
	}
}
```