## Expected Output

```text
skill-alpha  up to date
skill-beta  up to date

0 updated · 2 up to date · 0 not installed
```

## Expected

- Two column-0 status lines in registry order: `skill-alpha` then `skill-beta`, each `up to date`.
- Blank line then summary `0 updated · 2 up to date · 0 not installed`.
- Trailing newline; no leading blank line; no indented file lines; no legacy strings.

## Side Effects

- None (content already current).

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
skill-alpha  up to date
skill-beta  up to date

0 updated · 2 up to date · 0 not installed
`)
}
```
