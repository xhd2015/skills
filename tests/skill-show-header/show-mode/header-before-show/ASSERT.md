## Expected

- Command exits with code 0.
- stdout contains `---` and `name:`.
- stdout contains `name: go-best-practice`.
- stdout does **not** contain `# Go Best Practice Skill`.

## Side Effects

- None beyond process stdout.

## Errors

- No error from `Run`.

## Exit Code

- 0

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run failed: %v\nstderr:\n%s", err, resp.Stderr)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit code = %d, want 0\nstderr:\n%s", resp.ExitCode, resp.Stderr)
	}
	if !strings.Contains(resp.Stdout, "---") {
		t.Fatalf("stdout missing frontmatter delimiter:\n%s", resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "name:") {
		t.Fatalf("stdout missing name field:\n%s", resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "name: go-best-practice") {
		t.Fatalf("stdout missing skill name:\n%s", resp.Stdout)
	}
	if strings.Contains(resp.Stdout, "# Go Best Practice Skill") {
		t.Fatalf("stdout must not include body marker:\n%s", resp.Stdout)
	}
}
```