## Expected
- Exit code 0.
- stdout shows `State: closed`.
- stdout contains `#99`, `@dana`, `Merged feature`.
- stdout does NOT contain `#42` or `Still open`.

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
	for _, want := range []string{"State: closed", "#99", "@dana", "Merged feature"} {
		if !strings.Contains(stdout, want) {
			t.Fatalf("stdout missing %q:\n%s", want, stdout)
		}
	}
	if strings.Contains(stdout, "#42") || strings.Contains(stdout, "Still open") {
		t.Fatalf("stdout should not contain open PR:\n%s", stdout)
	}
}
```