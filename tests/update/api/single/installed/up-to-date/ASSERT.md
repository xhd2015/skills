## Expected Output

```text
skill-alpha  up to date
```

## Expected

- stdout is exactly `skill-alpha  up to date` plus trailing newline.
- No `updated` / file lines / legacy InstallTo strings.
- `SKILL.md` content remains canonical.

## Side Effects

- `SKILL.md` content remains canonical.

## Errors

- None.

```go
import (
	"path/filepath"
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
---
skill-alpha  up to date
`)
	skillPath := filepath.Join(resp.WorkDir, skillMDPath(skillAgentsDir("skill-alpha")))
	if got := readFile(t, skillPath); got != skillAlphaContent {
		t.Fatalf("SKILL.md content mismatch:\n%s", got)
	}
}
```
