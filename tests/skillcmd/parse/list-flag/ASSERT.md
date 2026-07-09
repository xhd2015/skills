## Expected

- Action is `list`.
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
	if resp.Action != "list" {
		t.Fatalf("Action = %q, want list", resp.Action)
	}
	if len(resp.Rest) != 0 {
		t.Fatalf("Rest = %v, want empty", resp.Rest)
	}
}
```
