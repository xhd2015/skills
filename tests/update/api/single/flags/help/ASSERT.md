## Expected Output

```text
<contains>
--dry-run
--global
</contains>
```

## Expected

- stdout documents `--dry-run` and `--global` (update flag surface).

## Side Effects

- No skill directories created.

## Errors

- None.

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
	for _, needle := range []string{"--dry-run", "--global"} {
		if !strings.Contains(resp.Stdout, needle) {
			t.Fatalf("stdout missing %q:\n%s", needle, resp.Stdout)
		}
	}
	if pathExists(t, filepath.Join(resp.WorkDir, skillAgentsDir("skill-alpha"))) {
		t.Fatalf("help must not create install dirs")
	}
}
```