## Expected

- Status line `skill-alpha  updated  (1 update)` with indented absolute path under the isolated `$HOME`.
- Global `SKILL.md` under `$HOME/.agents/skills/skill-alpha` is canonical.
- Project-local `.agents/skills/skill-alpha` does not exist.
- No legacy `Update skill at` header.

## Side Effects

- Only the global install location is updated.

## Errors

- None.

```go
import (
	"fmt"
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

	if resp.HomeDir == "" {
		t.Fatal("expected resp.HomeDir for global-scope leaf")
	}
	globalSkill := filepath.Join(resp.HomeDir, skillAgentsDir("skill-alpha"), "SKILL.md")
	// InstallTo uses filepath.Abs; on macOS that may prefix /private — accept either via regex.
	assert.Output(t, resp.Stdout, fmt.Sprintf(`---
version: 3
__PATH__: regex=.*/\.agents/skills/skill-alpha/SKILL\.md
---
skill-alpha  updated  \(1 update\)
  update  __PATH__
`))

	if got := readFile(t, globalSkill); got != skillAlphaContent {
		t.Fatalf("global SKILL.md not restored:\n%s", got)
	}
	localAgents := filepath.Join(resp.WorkDir, skillAgentsDir("skill-alpha"))
	if pathExists(t, localAgents) {
		t.Fatalf("local %s should not exist", localAgents)
	}
}
```
