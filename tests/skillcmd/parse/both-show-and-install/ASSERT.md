## Expected

- Returns an error (combined actions rejected).

## Side Effects

- None.

## Errors

- See Expected.

```go
import (
	"github.com/xhd2015/doctest/session"
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if resp.Error == "" {
		t.Fatalf("expected error for --show with --install, got Action=%q Rest=%v", resp.Action, resp.Rest)
	}
	// error should mention the conflicting actions
	if !strings.Contains(resp.Error, "show") && !strings.Contains(resp.Error, "install") {
		t.Fatalf("error should mention show/install: %s", resp.Error)
	}
}
```
