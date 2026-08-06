## Expected Output

```text
skill-alpha  updated  (1 update)
  update  <abs>/.codex/skills/skill-alpha/SKILL.md
```

## Expected

- One polished status line for the single skill (one installed target processed).
- Indented absolute codex `SKILL.md` path.
- `.opencode/skills/skill-alpha` does not exist.
- Codex `SKILL.md` restored to canonical content.

## Side Effects

- Codex `SKILL.md` restored to canonical content.

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
__SKILL__: regex=/.+/\.codex/skills/skill-alpha/SKILL\.md
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

	opencodeDir := filepath.Join(resp.WorkDir, skillOpencodeDir("skill-alpha"))
	if pathExists(t, opencodeDir) {
		t.Fatalf("expected %s not to exist", opencodeDir)
	}
	codexSkill := absUnder(resp.WorkDir, skillMDPath(skillCodexDir("skill-alpha")))
	if got := readFile(t, codexSkill); got != skillAlphaContent {
		t.Fatalf("codex SKILL.md not restored:\n%s", got)
	}
}
```
