## Expected

- Action is `install`.
- Rest includes `--global`.

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
	if resp.Action != "install" {
		t.Fatalf("Action = %q, want install", resp.Action)
	}
	if len(resp.Rest) != 1 || resp.Rest[0] != "--global" {
		t.Fatalf("Rest = %v, want [--global]", resp.Rest)
	}
}
```
