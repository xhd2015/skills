## Expected
- No error is returned.
- stdout contains `"Installed skill to:"`.
- The skill file exists at `.cursor/skills/example-skill/SKILL.md` within the working directory.
- The skill file content is `"# test skill\n"`.

```go
import (
	"os"
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
	if !strings.Contains(resp.Stdout, "Installed skill to:") {
		t.Fatalf("stdout missing 'Installed skill to:':\n%s", resp.Stdout)
	}

	skillFile := filepath.Join(resp.WorkDir, ".cursor", "skills", "example-skill", "SKILL.md")
	content, readErr := os.ReadFile(skillFile)
	if readErr != nil {
		t.Fatalf("read skill file %s: %v", skillFile, readErr)
	}
	if string(content) != "# test skill\n" {
		t.Fatalf("unexpected skill content at %s: %q", skillFile, string(content))
	}
}
```
