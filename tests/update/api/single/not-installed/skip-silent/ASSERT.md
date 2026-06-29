## Expected

- `Run` succeeds with no error.
- stdout is empty (silent skip).

## Side Effects

- `.agents/skills/skill-alpha` is not created.

## Errors

- None.

## Exit Code

- N/A (library call).

```go
import (
	"path/filepath"
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
	if strings.TrimSpace(resp.Stdout) != "" {
		t.Fatalf("expected empty stdout, got:\n%s", resp.Stdout)
	}
	agents := filepath.Join(resp.WorkDir, skillAgentsDir("skill-alpha"))
	if pathExists(t, agents) {
		t.Fatalf("expected %s not to exist", agents)
	}
}
```