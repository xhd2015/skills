## Expected Output

```text
skill-alpha  updated  (1 create, 1 update)
  create  <abs>/.agents/skills/skill-alpha/extra.md
  update  <abs>/.agents/skills/skill-alpha/SKILL.md
skill-beta  up to date

1 updated · 1 up to date · 0 not installed
```

## Expected

- Alpha status: `updated` with `(1 create, 1 update)` counts.
- Indented file lines (2 spaces), op then two spaces then **absolute** paths; create before update.
- Beta: `up to date` with no file lines.
- Summary: `1 updated · 1 up to date · 0 not installed`.
- On-disk alpha `SKILL.md` restored; `extra.md` created.

## Side Effects

- Drifted content replaced; extra file written.

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
__EXTRA__: regex=/.+/extra\.md
__SKILL__: regex=/.+/SKILL\.md
---
skill-alpha  updated  \(1 create, 1 update\)
  create  __EXTRA__
  update  __SKILL__
skill-beta  up to date

1 updated · 1 up to date · 0 not installed
`)

	// Paths must be absolute and under the workdir skill tree.
	if !strings.Contains(resp.Stdout, "create  ") || !strings.Contains(resp.Stdout, "update  ") {
		t.Fatalf("expected create/update file lines with two-space separators:\n%s", resp.Stdout)
	}
	for _, line := range strings.Split(resp.Stdout, "\n") {
		trim := strings.TrimSpace(line)
		if strings.HasPrefix(trim, "create  ") || strings.HasPrefix(trim, "update  ") {
			path := strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(trim, "create  "), "update  "))
			if !filepath.IsAbs(path) {
				t.Fatalf("file path must be absolute: %q", path)
			}
			if !strings.Contains(path, skillAgentsDir("skill-alpha")) {
				t.Fatalf("file path not under skill-alpha tree: %q", path)
			}
		}
	}

	skillMD := absUnder(resp.WorkDir, skillMDPath(skillAgentsDir("skill-alpha")))
	extraMD := absUnder(resp.WorkDir, filepath.Join(skillAgentsDir("skill-alpha"), "extra.md"))
	if got := readFile(t, skillMD); got != skillAlphaContent {
		t.Fatalf("alpha SKILL.md not restored:\n%s", got)
	}
	if got := readFile(t, extraMD); got != extraFileContent {
		t.Fatalf("extra.md content mismatch:\n%s", got)
	}
}
```
