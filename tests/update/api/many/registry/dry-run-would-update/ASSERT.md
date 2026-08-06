## Expected Output

```text
skill-alpha  would update  (1 update)
  update  <abs>/.agents/skills/skill-alpha/SKILL.md
skill-beta  up to date

0 updated · 1 would update · 1 up to date · 0 not installed  [dry-run]
```

## Expected

- Alpha: `would update` with indented planned absolute path (no write).
- Beta: `up to date`.
- Summary includes `would update` and trailing `[dry-run]` meta.
- No per-line `[dry-run]` InstallTo prefixes; no legacy strings.
- Alpha `SKILL.md` remains drifted.

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
skill-beta  up to date

0 updated · 1 would update · 1 up to date · 0 not installed  \[dry-run\]
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

	skillMD := absUnder(resp.WorkDir, skillMDPath(skillAgentsDir("skill-alpha")))
	wantDrift := "# drifted for dry-run\n"
	if got := readFile(t, skillMD); got != wantDrift {
		t.Fatalf("SKILL.md should remain drifted, got:\n%s", got)
	}
}
```
