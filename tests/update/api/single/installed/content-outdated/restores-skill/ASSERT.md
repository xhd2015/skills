## Expected Output

```text
skill-alpha  updated  (1 update)
  update  <abs>/.agents/skills/skill-alpha/SKILL.md
```

## Expected

- Status line `skill-alpha  updated  (1 update)` at column 0.
- Indented absolute `update` file line (two spaces, no colon).
- `SKILL.md` on disk equals canonical embedded content.
- No legacy `Update skill at` header.

## Side Effects

- Drifted content is replaced.

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
skill-alpha  updated  \(1 update\)
  update  __SKILL__
`)

	for _, line := range strings.Split(resp.Stdout, "\n") {
		trim := strings.TrimSpace(line)
		if strings.HasPrefix(trim, "update  ") {
			path := strings.TrimSpace(strings.TrimPrefix(trim, "update  "))
			if !filepath.IsAbs(path) {
				t.Fatalf("file path must be absolute: %q", path)
			}
		}
	}

	skillPath := absUnder(resp.WorkDir, skillMDPath(skillAgentsDir("skill-alpha")))
	if got := readFile(t, skillPath); got != skillAlphaContent {
		t.Fatalf("SKILL.md not restored, got:\n%s", got)
	}
}
```
