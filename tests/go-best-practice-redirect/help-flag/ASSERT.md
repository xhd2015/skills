## Expected

- Exit code `1`.
- stderr contains `https://github.com/xhd2015/go-best-practice` and module + `go install` path (shared redirect).
- stdout is empty — no full old help (Commands / vet / topics listing).

## Side Effects

- None beyond process I/O.

## Errors

- No error from `Run`.

## Exit Code

- 1

```go
import (
	"strings"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertRedirect(t, resp, err)
	// Extra guard: old full CLI --help listed Commands and vet on stdout.
	if strings.Contains(resp.Stdout, "vet [flags]") {
		t.Fatalf("stdout must not contain old full help (vet):\n%s", resp.Stdout)
	}
	if strings.Contains(resp.Stdout, "Usage: go-best-practice <command>") {
		t.Fatalf("stdout must not contain old full usage:\n%s", resp.Stdout)
	}
}
```
