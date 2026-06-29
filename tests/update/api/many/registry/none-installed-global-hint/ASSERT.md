## Expected

- stdout contains `skill not installed: skill-alpha` and `skill not installed: skill-beta`.
- stdout does not contain `No installed skills found`.

## Expected Output

```
<contains>
skill not installed: skill-alpha
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
	assert.Output(t, resp.Stdout, `
<contains>
skill not installed: skill-alpha
skill not installed: skill-beta
</contains>`)
	if strings.Contains(resp.Stdout, "No installed skills found") {
		t.Fatalf("aggregate scope hint must be removed:\n%s", resp.Stdout)
	}
}
```