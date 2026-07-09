## Expected

- Returns an error (no --show/--install/--list).

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
	if resp.Error == "" {
		t.Fatalf("expected error for missing action, got Action=%q Rest=%v", resp.Action, resp.Rest)
	}
}
```
