## Expected

- Exit code `1`.
- stderr contains `https://github.com/xhd2015/go-best-practice`.
- stderr contains `github.com/xhd2015/go-best-practice`.
- stderr contains `go install github.com/xhd2015/go-best-practice/cmd/go-best-practice@latest`.
- stdout is empty (no topic list / product output).

## Side Effects

- None beyond process I/O.

## Errors

- No error from `Run` (non-zero process exit is captured, not returned).

## Exit Code

- 1

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertRedirect(t, resp, err)
}
```
