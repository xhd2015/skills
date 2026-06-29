## Expected

- stdout contains `[dry-run]`.
- stdout does not contain a non-dry-run `Update skill at` line (only dry-run prefixed output).
- `SKILL.md` still contains drifted content.

## Side Effects

- No write of canonical content.

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
	if !strings.Contains(resp.Stdout, "[dry-run]") {
		t.Fatalf("stdout missing dry-run prefix:\n%s", resp.Stdout)
	}
	if strings.Contains(resp.Stdout, "Update skill at") {
		t.Fatalf("stdout must not contain non-dry-run update line:\n%s", resp.Stdout)
	}
	skillPath := filepath.Join(resp.WorkDir, skillMDPath(skillAgentsDir("skill-alpha")))
	wantDrift := "# drifted for dry-run\n"
	if got := readFile(t, skillPath); got != wantDrift {
		t.Fatalf("SKILL.md should remain drifted, got:\n%s", got)
	}
}
```