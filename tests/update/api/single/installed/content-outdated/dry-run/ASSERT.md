## Expected Output

```text
skill-alpha  would update  (1 update)
  update  <abs>/.agents/skills/skill-alpha/SKILL.md
```

## Expected

- Status is `would update` (not applied `updated`).
- Indented planned absolute path; no legacy `[dry-run] Update skill at` headers.
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

	"github.com/xhd2015/doctest/assert"
	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if resp.Error != "" {
		t.Fatalf("expected no error, got: %s", resp.Error)
	}
	assertBatchStdoutPolished(t, resp.Stdout)
	assertNoLegacyUpdateStdout(t, resp.Stdout)

	assert.Output(t, resp.Stdout, `---
version: 3
__SKILL__: regex=/.+/SKILL\.md
---
skill-alpha  would update  \(1 update\)
  update  __SKILL__
`)

	for _, line := range strings.Split(resp.Stdout, "\n") {
		trim := strings.TrimSpace(line)
		if strings.HasPrefix(trim, "update  ") {
			path := strings.TrimSpace(strings.TrimPrefix(trim, "update  "))
			if !filepath.IsAbs(path) {
				t.Fatalf("planned path must be absolute: %q", path)
			}
		}
	}

	skillPath := absUnder(resp.WorkDir, skillMDPath(skillAgentsDir("skill-alpha")))
	wantDrift := "# drifted for dry-run\n"
	if got := readFile(t, skillPath); got != wantDrift {
		t.Fatalf("SKILL.md should remain drifted, got:\n%s", got)
	}
}
```
