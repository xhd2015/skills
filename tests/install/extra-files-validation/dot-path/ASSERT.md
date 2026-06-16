## Expected
- An error is returned (stored in `resp.Error`).
- The error message contains `"invalid install file path"`.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if resp.Error == "" {
		t.Fatal("expected an error for invalid extra file path '.'")
	}
	if !strings.Contains(resp.Error, "invalid install file path") {
		t.Fatalf("expected 'invalid install file path' in error, got: %s", resp.Error)
	}
}
```
