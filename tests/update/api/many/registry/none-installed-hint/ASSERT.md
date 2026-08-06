## Expected Output

```text
skill-alpha  not installed
skill-beta  not installed

0 updated · 0 up to date · 2 not installed
```

## Expected

- `Run` succeeds with no error.
- Each registry skill prints `name  not installed` (two spaces; no `skill not installed:` prefix).
- Summary counts two not-installed; no legacy aggregate hint.
- No skill directories created under the workdir.

## Side Effects

- No skill directories created under the workdir.

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
skill-alpha  not installed
skill-beta  not installed

0 updated · 0 up to date · 2 not installed
`)
	for _, name := range []string{"skill-alpha", "skill-beta"} {
		if pathExists(t, filepath.Join(resp.WorkDir, skillAgentsDir(name))) {
			t.Fatalf("expected no install dir for %s", name)
		}
	}
}
```
