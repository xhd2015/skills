## Expected Output

```text
skill-alpha  up to date
skill-beta  not installed

0 updated · 1 up to date · 1 not installed
```

## Expected

- Alpha is up to date; beta is not installed.
- `skill-beta` agents directory does not exist after the run.
- No legacy strings; polished summary.

## Side Effects

- No install created for beta.

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
skill-beta  not installed

0 updated · 1 up to date · 1 not installed
`)
	if pathExists(t, filepath.Join(resp.WorkDir, skillAgentsDir("skill-beta"))) {
		t.Fatalf("beta install dir must not be created")
	}
}
```
