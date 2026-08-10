# Playwright-Debug NODE_PATH + Explicit Page Param Tests

Doc-style tests for playwright-debug fixes that (A) set NODE_PATH to the
playwright cache node_modules when launching node for file/eval mode, and
(B) require nested library functions to receive `page` explicitly instead of
relying on an implicit global.

# DSN (Domain Specific Notion)

Participants:

- User invokes playwright-debug `run <file.js>` to execute a script with
  injected browser, page, chromium, and require.
- CLI file runner ensures playwright is cached under
  `~/.playwright-debug/node_package/`, writes bootstrap.cjs, and spawns node
  with environment including NODE_PATH pointing at the cache node_modules.
- Nested user modules loaded via require() run in isolated CommonJS scope; they
  must receive `page` as an explicit parameter — no implicit global.
- Nested modules that require('playwright') resolve via NODE_PATH fallback to
  the cached playwright package when not installed locally.

## Decision Tree

```text
playwright-debug-node-path-page/
├── node-path/
│   ├── nested-require-playwright/
│   └── env-not-empty/
└── explicit-page/
    └── nested-lib-uses-page/
```

## Test Index

- node-path/nested-require-playwright: nested lib require('playwright') via NODE_PATH → playwright-ok (slow)
- node-path/env-not-empty: NODE_PATH env contains playwright-debug cache path → node-path-set (slow)
- explicit-page/nested-lib-uses-page: nested lib uses explicit page param → explicit-page-ok (slow)

## How to Run

```sh
doctest vet ./tests/playwright-debug-node-path-page
doctest test -v ./tests/playwright-debug-node-path-page
doctest test -v --label slow ./tests/playwright-debug-node-path-page
```

## Version

0.0.2

```go
import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/xhd2015/doctest/session"
)

type Request struct {
	Args []string
}

type Response struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	_ = d
	bin, err := buildPlaywrightDebugOnce(t, d)
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(bin, req.Args...)
	cmd.Env = append(os.Environ(), "GOWORK=off")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	runErr := cmd.Run()
	resp := &Response{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}
	if cmd.ProcessState != nil {
		resp.ExitCode = cmd.ProcessState.ExitCode()
	}
	if runErr != nil {
		if _, ok := runErr.(*exec.ExitError); !ok {
			return resp, fmt.Errorf("run playwright-debug: %w", runErr)
		}
	}
	return resp, nil
}
```