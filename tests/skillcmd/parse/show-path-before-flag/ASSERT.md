## Expected

- Action is `show`.
- Rest is `[flags-parsing/types]`.

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
	if len(resp.Rest) != 1 || resp.Rest[0] != "flags-parsing/types" {
		t.Fatalf("Rest = %v, want [flags-parsing/types]", resp.Rest)
	}
}
```
