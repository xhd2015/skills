## Expected

- stdout contains exactly one `Skill is up to date` line referencing `skill-alpha`.
- stdout contains `skill not installed: skill-beta`.
- Lines appear in registry CLI-name order: alpha status before beta not-installed.
- stdout does not contain `No installed skills found`.

## Expected Output

```
<contains>
Skill is up to date
skill-alpha
skill not installed: skill-beta
</contains>
```

## Errors

- None.

```go
import (
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
	assert.Output(t, resp.Stdout, `` +
`<contains>
Skill is up to date
skill-alpha
skill not installed: skill-beta
</contains>`)
	if strings.Count(resp.Stdout, "Skill is up to date") != 1 {
		t.Fatalf("expected one up-to-date line:\n%s", resp.Stdout)
	}
	if strings.Contains(resp.Stdout, "No installed skills found") {
		t.Fatalf("aggregate scope hint must be removed:\n%s", resp.Stdout)
	}
	alphaIdx := strings.Index(resp.Stdout, "skill-alpha")
	betaIdx := strings.Index(resp.Stdout, "skill not installed: skill-beta")
	if alphaIdx < 0 || betaIdx < 0 || alphaIdx > betaIdx {
		t.Fatalf("expected alpha output before beta not-installed line:\n%s", resp.Stdout)
	}
}
```
