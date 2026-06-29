## Expected

- stdout contains two lines with `Skill is up to date`.
- stdout references both `skill-alpha` and `skill-beta` install paths.

## Side Effects

- None (content already current).

## Errors

- None.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if resp.Error != "" {
		t.Fatalf("expected no error, got: %s", resp.Error)
	}
	count := strings.Count(resp.Stdout, "Skill is up to date")
	if count != 2 {
		t.Fatalf("expected 2 up-to-date lines, got %d:\n%s", count, resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "skill-alpha") {
		t.Fatalf("stdout missing alpha:\n%s", resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "skill-beta") {
		t.Fatalf("stdout missing beta:\n%s", resp.Stdout)
	}
}
```