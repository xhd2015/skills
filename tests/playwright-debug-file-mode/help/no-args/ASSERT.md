## Expected

- Exit code 0.
- stdout contains `Usage:`.
- stdout mentions `run` with a `.js` file argument.
- stdout mentions `-e` or `--eval` for adhoc eval.
- stdout indicates bare existing `.js` path is accepted (file alias).
- stdout states `run` requires an existing script file.
- stdout mentions trailing script arguments are passed through to `process.argv`.
- stdout ends with a newline.

## Side Effects

- None beyond printing help.

## Errors

- No error from `Run`.

## Exit Code

- 0

```go
import (
	"strings"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run failed: %v\nstderr:\n%s", err, resp.Stderr)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit code = %d, want 0\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	out := resp.Stdout
	if !strings.Contains(out, "Usage:") {
		t.Fatalf("stdout missing Usage:\n%s", out)
	}
	lower := strings.ToLower(out)
	if !strings.Contains(lower, "run") {
		t.Fatalf("stdout must mention run subcommand:\n%s", out)
	}
	if !strings.Contains(lower, ".js") {
		t.Fatalf("stdout must mention .js file:\n%s", out)
	}
	if !strings.Contains(out, "-e") && !strings.Contains(out, "--eval") {
		t.Fatalf("stdout must mention -e or --eval:\n%s", out)
	}
	if !strings.Contains(lower, "file") {
		t.Fatalf("stdout must describe file mode or file alias:\n%s", out)
	}
	if !strings.Contains(out, "process.argv") {
		t.Fatalf("stdout must mention process.argv for script arg pass-through:\n%s", out)
	}
	if len(out) == 0 || out[len(out)-1] != '\n' {
		t.Fatalf("stdout must end with trailing newline:\n%q", out)
	}
}
```