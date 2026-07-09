## Expected

- Action is `show`.
- Header is true.
- Rest is empty.

## Side Effects

- None.

## Errors

- See Expected.

```go
import (
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if resp.Error != "" {
		t.Fatalf("unexpected error: %s", resp.Error)
	}
	if resp.Action != "show" {
		t.Fatalf("Action = %q, want show", resp.Action)
	}
	if !resp.Header {
		t.Fatalf("Header = false, want true")
	}
	if len(resp.Rest) != 0 {
		t.Fatalf("Rest = %v, want empty", resp.Rest)
	}
}
```
