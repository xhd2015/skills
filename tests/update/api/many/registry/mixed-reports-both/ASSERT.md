## Expected Output

```text
skill-alpha  up to date
skill-beta  not installed

0 updated · 1 up to date · 1 not installed
```

## Expected

- Lines appear in registry CLI-name order: alpha status before beta not-installed.
- Summary reflects one up-to-date and one not-installed.
- No legacy `Skill is up to date:` / `skill not installed:` forms.

## Errors

- None.

```go
import (
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
---
skill-alpha  up to date
skill-beta  not installed

0 updated · 1 up to date · 1 not installed
`)
	alphaIdx := strings.Index(resp.Stdout, "skill-alpha  up to date")
	betaIdx := strings.Index(resp.Stdout, "skill-beta  not installed")
	if alphaIdx < 0 || betaIdx < 0 || alphaIdx > betaIdx {
		t.Fatalf("expected alpha status before beta not-installed:\n%s", resp.Stdout)
	}
}
```
