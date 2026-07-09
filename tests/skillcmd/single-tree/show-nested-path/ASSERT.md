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
	if resp.Error != "" {
		t.Fatalf("unexpected error: %s", resp.Error)
	}
	if !strings.Contains(resp.Stdout, "# Nested A/B Body") {
		t.Fatalf("stdout missing nested body:\n%s", resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "name: demo-skill/a/b") {
		t.Fatalf("stdout missing nested name:\n%s", resp.Stdout)
	}
	if strings.Contains(resp.Stdout, "# Demo Skill Body") {
		t.Fatalf("must not print root body for nested path:\n%s", resp.Stdout)
	}
}
```
