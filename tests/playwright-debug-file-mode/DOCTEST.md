# Playwright-Debug File Mode + Eval Flag Tests

Doc-style tests for cmd/playwright-debug CLI routing between file mode
(existing .js scripts with require and top-level await) and eval mode
(adhoc snippets via -e / --eval or bare non-file args).

# DSN (Domain Specific Notion)

Participants:

- User invokes the playwright-debug CLI with positional args, run, or
  -e / --eval flags.
- CLI router (handle) classifies input into help, explicit file mode
  (run with file.js or bare existing .js path), or eval mode (-e, --eval,
  or bare non-file script). Routing errors must fail before browser startup.
- File runner writes embedded bootstrap.cjs to the playwright cache dir,
  chdirs to the script directory, and executes the user file via AsyncFunction
  with injected browser, page, chromium, require, __filename, and __dirname.
- Eval runner wraps the snippet in an async IIFE and runs node -e (existing
  behavior).
- Playwright cache (~/.playwright-debug/node_package/) supplies the
  playwright npm package and Chromium for script execution leaves.

## Decision Tree

```text
playwright-debug-file-mode/
├── help/
│   ├── no-args/
│   └── with-flag/
├── run-command/
│   ├── errors/
│   │   ├── missing-arg/
│   │   ├── non-js-script/
│   │   ├── file-not-found/
│   │   └── too-many-args/
│   └── execution/
│       ├── simple-file/
│       ├── toplevel-await/
│       └── with-require/
├── eval-flag/
│   ├── errors/
│   │   ├── missing-script/
│   │   └── extra-args/
│   └── execution/
│       └── short-flag/
└── bare-arg/
    ├── file-alias/
    │   └── simple-file/
    └── eval-script/
        └── console-log/
```

## Test Index

- help/no-args: no args prints usage, exit 0 (fast)
- help/with-flag: -h prints usage, exit 0 (fast)
- run-command/errors/missing-arg: run alone → file required error (fast)
- run-command/errors/non-js-script: run with inline await snippet → requires existing .js file (fast)
- run-command/errors/file-not-found: run missing.js → file not found (fast)
- run-command/errors/too-many-args: run with two files → exactly one file (fast)
- eval-flag/errors/missing-script: -e without script → error (fast)
- eval-flag/errors/extra-args: --eval with extra positional arg → error (fast)
- run-command/execution/simple-file: run simple-eval.js → file-mode-ok (slow)
- run-command/execution/toplevel-await: top-level await on about:blank (slow)
- run-command/execution/with-require: relative require() resolves (slow)
- eval-flag/execution/short-flag: -e console.log eval-ok (slow)
- bare-arg/file-alias/simple-file: bare .js path routes to file mode (slow)
- bare-arg/eval-script/console-log: bare script string routes to eval mode (slow)

## How to Run

```sh
doctest vet ./tests/playwright-debug-file-mode
doctest test -v ./tests/playwright-debug-file-mode
doctest test -v --label slow ./tests/playwright-debug-file-mode
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
)

type Request struct {
	Args []string
}

type Response struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

func Run(t *testing.T, req *Request) (*Response, error) {
	bin, err := buildPlaywrightDebugOnce(t)
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