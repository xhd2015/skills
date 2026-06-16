## Expected
- Exit code 0 with same content as `open-default` (issues only, PR excluded).

## Exit Code
- 0

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected exit code 0, got %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	stdout := resp.Stdout
	for _, want := range []string{
		"testowner/testrepo",
		"#15", "Login page broken",
		"#14", "Add dark mode",
		"Page 1 (2 items)",
	} {
		if !strings.Contains(stdout, want) {
			t.Fatalf("stdout missing %q:\n%s", want, stdout)
		}
	}
	if strings.Contains(stdout, "#42") {
		t.Fatalf("stdout should not contain PR #42:\n%s", stdout)
	}
}
```