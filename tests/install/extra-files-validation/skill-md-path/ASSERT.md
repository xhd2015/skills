## Expected
- An error is returned (stored in `resp.Error`).
- The error message contains `"extra install file cannot replace SKILL.md"`.

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
		t.Fatal("expected an error for extra file path 'SKILL.md'")
	}
	if !strings.Contains(resp.Error, "extra install file cannot replace SKILL.md") {
		t.Fatalf("expected 'extra install file cannot replace SKILL.md' in error, got: %s", resp.Error)
	}
}
```
