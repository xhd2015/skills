## Expected

- Scenario-specific success or error as asserted.

## Side Effects

- install-extra-files creates nested skill-cli/SKILL.md under the skill dir.

## Errors

- reject-dotdot expects a non-empty error.

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
		t.Fatalf("expected error for ../x, stdout:\n%s", resp.Stdout)
	}
	// accept either invalid segment or unknown topic wording
	errLow := strings.ToLower(resp.Error)
	if !strings.Contains(errLow, "..") && !strings.Contains(errLow, "invalid") && !strings.Contains(errLow, "path") {
		t.Fatalf("error should mention invalid path: %s", resp.Error)
	}
}
```
