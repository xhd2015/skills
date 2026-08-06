## Expected

- Action is `show`.
- Header is false.
- Rest is empty.
- No error.

## Side Effects

- None.

## Errors

- See Expected.

```go
import (
	"github.com/xhd2015/doctest/session"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if resp.Error != "" {
		t.Fatalf("unexpected error: %s", resp.Error)
	}
	if resp.Action != "show" {
		t.Fatalf("Action = %q, want show", resp.Action)
	}
	if resp.Header {
		t.Fatalf("Header = true, want false")
	}
	if len(resp.Rest) != 0 {
		t.Fatalf("Rest = %v, want empty", resp.Rest)
	}
}
```
