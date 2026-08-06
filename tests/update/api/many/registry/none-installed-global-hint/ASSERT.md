## Expected Output

```text
skill-alpha  not installed
skill-beta  not installed

0 updated · 0 up to date · 2 not installed
```

## Expected

- Same not-installed shape as local none-installed (no aggregate scope hint).
- Summary with two not-installed; no legacy `skill not installed:` prefix.

## Errors

- None.

```go
import (
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
}
```
