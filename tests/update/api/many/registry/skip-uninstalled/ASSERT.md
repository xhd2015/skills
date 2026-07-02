## Expected

- stdout contains `Skill is up to date` for `skill-alpha`.
- stdout contains `skill not installed: skill-beta`.
- `skill-beta` agents directory does not exist.

## Side Effects

- No install created for beta.

## Errors

- None.

```go
import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/assert"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if resp.Error != "" {
		t.Fatalf("expected no error, got: %s", resp.Error)
	}
	assertBatchStdoutPolished(t, resp.Stdout)
	if !strings.Contains(resp.Stdout, "Skill is up to date") {
		t.Fatalf("expected up-to-date line for alpha:\n%s", resp.Stdout)
	}
	assert.Output(t, resp.Stdout, `` +
`<contains>
skill not installed: skill-beta
</contains>`)
	if pathExists(t, filepath.Join(resp.WorkDir, skillAgentsDir("skill-beta"))) {
		t.Fatalf("beta install dir must not be created")
	}
}
```
