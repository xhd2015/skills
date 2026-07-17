## Expected

- No error.
- stdout lists the skill name first, then topic paths derived from `**/TOPIC.md`
  (slash paths without the filename), sorted:
  - `demo-skill`
  - `a/b`
  - `skill-cli`
- Trailing newline after last line.

## Side Effects

- None.

## Errors

- None.

```go
import (
	"fmt"
	"testing"

	"github.com/xhd2015/doctest/assert"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if resp.Error != "" {
		t.Fatalf("unexpected error: %s", resp.Error)
	}
	assert.Output(t, resp.Stdout, fmt.Sprintf(`---
version: 3
---
%s
a/b
skill-cli
`, demoSkillName))
}
```
