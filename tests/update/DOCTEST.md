# Update Package Tests

Doc-style tests for `install.HandleUpdate` and `install.HandleUpdateMany`.
Update refreshes **only** directories where `SKILL.md` already exists; missing
installs are never created. Batch `HandleUpdateMany` prints a polished per-skill
status line (and optional indented file ops) for every registry entry, then a
trailing summary. Single `HandleUpdate` stays silent when nothing is installed
and uses the same status/file-line dialect when installed.

# DSN (Domain Specific Notion)

Participants:

- **Update handler** — parses update-only flags (`--global`, target tool flags,
  `--dry-run`, `--help`), resolves target directories, and never calls
  `HandleInstall`.
- **Target resolver** — maps flags to the same directory layout as install
  (`.agents/skills/<name>`, `.cursor/skills/<name>`, etc., optionally under
  `$HOME` when `--global` is set).
- **Install probe** — treats a target as installed when `<dir>/SKILL.md` exists;
  otherwise the skill is reported `not installed` (batch) or skipped silently
  (single).
- **Inventory apply** — when a target is installed, compares embedded skill
  content (and extra files) to disk; writes create/update/delete when not
  dry-run; emits polished status + indented absolute file paths.
- **Batch reporter** — `HandleUpdateMany` walks registry order, prints one
  status line per skill, optional 2-space-indented file lines, then a trailing
  count summary (and `[dry-run]` meta when applicable).
- **Isolated runner** — test harness process (`cmd/runupdate`) chdirs / sets
  `HOME` / owns product stdout so suite `Run` stays Parallel-safe.

Behaviors:

- **Single skill** — `HandleUpdate(opts, args)` updates one skill definition
  across all resolved targets that are installed; silent when none installed.
- **Registry batch** — `HandleUpdateMany(skills, args)` walks each registry
  entry in stable CLI-name order; always emits status + trailing summary.

## Decision Tree

```text
update/
└── api/                              # Entry API (largest behavioral split)
    ├── single/                       # HandleUpdate
    │   ├── not-installed/
    │   │   └── skip-silent
    │   ├── installed/
    │   │   ├── up-to-date
    │   │   └── content-outdated/
    │   │       ├── restores-skill
    │   │       └── dry-run
    │   ├── target-flags/
    │   │   └── partial-dirs
    │   ├── global-scope/
    │   │   └── global-home
    │   └── flags/
    │       ├── help
    │       └── color-conflict
    └── many/                         # HandleUpdateMany
        └── registry/
            ├── none-installed-hint
            ├── none-installed-global-hint
            ├── skip-uninstalled
            ├── mixed-reports-both
            ├── all-up-to-date
            ├── updated-with-files
            ├── dry-run-would-update
            ├── color-force-on
            └── color-force-off
```

## Test Index

| Leaf | Description |
|------|-------------|
| `api/single/not-installed/skip-silent` | No `SKILL.md` → empty stdout, no directories created |
| `api/single/installed/up-to-date` | Installed + matching content → `name  up to date` |
| `api/single/installed/content-outdated/restores-skill` | Drifted `SKILL.md` → `updated` + indented `update` path; content restored |
| `api/single/installed/content-outdated/dry-run` | Drift + `--dry-run` → `would update` + paths; file unchanged |
| `api/single/target-flags/partial-dirs` | Installed under `--codex` only; `--codex --opencode` update touches codex only |
| `api/single/global-scope/global-home` | `--global` install + update; global tree updated, local default untouched |
| `api/single/flags/help` | `--help` prints update usage (incl. `--color` / `--no-color`; not `--update` flag) |
| `api/single/flags/color-conflict` | `--color` + `--no-color` → mutual exclusion error |
| `api/many/registry/none-installed-hint` | No installs → `name  not installed` per skill + summary |
| `api/many/registry/none-installed-global-hint` | `--global`, no installs → same not-installed shape |
| `api/many/registry/skip-uninstalled` | One installed + one missing → up-to-date + not-installed; beta dir absent |
| `api/many/registry/mixed-reports-both` | Mixed batch in registry order + summary counts |
| `api/many/registry/all-up-to-date` | Two installed, current → two `up to date` + summary |
| `api/many/registry/updated-with-files` | Drift + extra planned file → `updated` + indented create/update abs paths |
| `api/many/registry/dry-run-would-update` | Drift + `--dry-run` → `would update` + paths; disk unchanged; summary `[dry-run]` |
| `api/many/registry/color-force-on` | `--color` forces ANSI (gray SGR) on pipe |
| `api/many/registry/color-force-off` | `--no-color` emits no ANSI |

## Output contract (batch)

Status lines start at column 0: `{name}  {status}` (exactly two spaces before status).

| Status | When | File lines |
|--------|------|------------|
| `up to date` | installed, no content change | none |
| `updated` | installed, files written | yes, 2-space indent + abs paths |
| `would update` | `--dry-run` and would write | yes, planned abs paths |
| `not installed` | no `SKILL.md` at resolved targets | none |

When there are file ops, status carries counts:

```text
followup  updated  (1 create, 6 update)
  create  /abs/.../host/TOPIC.md
  update  /abs/.../SKILL.md
```

File lines: two-space indent, op word, two spaces, absolute path (no colon after op).
Op order: create, then update, then delete. Only non-zero ops appear in the
parenthetical, same order.

Trailing summary (blank line before it; stdout ends with `\n`; no leading blank):

```text
N updated · M up to date · K not installed
```

Dry-run:

```text
0 updated · W would update · M up to date · K not installed  [dry-run]
```

Always include `updated`, `up to date`, and `not installed` counts (zeros OK).
Include `would update` only under `--dry-run`.

**Single-skill dialect:** same status + file-line shape; no multi-count summary
required (batch-only). Silent when not installed.

**Legacy forbidden** in status/batch output: `Skill is up to date:`,
`skill not installed:`, bare trailing skill-name-only lines after update,
`Update skill at ` path headers from InstallTo.

## How to Run

```sh
doctest vet ./tests/update
doctest test -v ./tests/update
```

Classic TDD: leaves describe the **new** format and are expected **RED** until
`skillcmd` implements the polished printer.

## Version

0.0.2

```go
import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/session"
	"github.com/xhd2015/skills/install"
)

type FileMutate struct {
	RelPath string `json:"relPath"`
	Content string `json:"content"`
}

type PreInstall struct {
	Opts install.InstallOptions `json:"opts"`
	Args []string               `json:"args"`
}

type Request struct {
	UseMany bool

	SingleOpts install.InstallOptions
	ManySkills []install.UpdateSkill

	UpdateArgs []string

	PreInstalls       []PreInstall
	PostInstallMutate []FileMutate

	UseGlobalHome bool
}

type Response struct {
	Stdout  string
	Error   string
	WorkDir string
	// HomeDir is the isolated HOME used for --global leaves (empty otherwise).
	HomeDir string
}

// harnessRequest is the JSON payload for cmd/runupdate (isolated process).
type harnessRequest struct {
	UseMany           bool                   `json:"useMany"`
	SingleOpts        install.InstallOptions `json:"singleOpts"`
	ManySkills        []install.UpdateSkill  `json:"manySkills"`
	UpdateArgs        []string               `json:"updateArgs"`
	PreInstalls       []PreInstall           `json:"preInstalls"`
	PostInstallMutate []FileMutate           `json:"postInstallMutate"`
	UseGlobalHome     bool                   `json:"useGlobalHome"`
	WorkDir           string                 `json:"workDir"`
	HomeDir           string                 `json:"homeDir"`
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	tmpDir := t.TempDir()
	homeDir := ""
	if req.UseGlobalHome {
		homeDir = t.TempDir()
	}

	hr := harnessRequest{
		UseMany:           req.UseMany,
		SingleOpts:        req.SingleOpts,
		ManySkills:        req.ManySkills,
		UpdateArgs:        req.UpdateArgs,
		PreInstalls:       req.PreInstalls,
		PostInstallMutate: req.PostInstallMutate,
		UseGlobalHome:     req.UseGlobalHome,
		WorkDir:           tmpDir,
		HomeDir:           homeDir,
	}
	reqFile := filepath.Join(tmpDir, "runupdate-req.json")
	raw, err := json.Marshal(hr)
	if err != nil {
		return nil, err
	}
	if err := os.WriteFile(reqFile, raw, 0644); err != nil {
		return nil, err
	}

	// d.DOCTEST_ROOT is tests/update — harness lives at cmd/runupdate.
	// Child process owns chdir / HOME / product stdout (Parallel-safe parent).
	harnessDir := filepath.Join(d.DOCTEST_ROOT, "cmd", "runupdate")
	// Args after the package path are program args (not go-run flags).
	cmd := exec.Command("go", "run", ".", "-req", reqFile)
	cmd.Dir = harnessDir
	// Isolate from parent go.work (consumer agent-pro workspace) so the skills
	// module is the main module for the harness package.
	cmd.Env = append(os.Environ(), "GOWORK=off")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	runErr := cmd.Run()

	resp := &Response{
		Stdout:  stdout.String(),
		WorkDir: tmpDir,
		HomeDir: homeDir,
	}
	if runErr != nil {
		// Product / harness failure: stderr carries the error string from child.
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = runErr.Error()
		}
		resp.Error = msg
		// Harness isolation errors (exit 2) surface as Run error; product errors
		// (exit 1) stay on resp.Error for Assert to inspect.
		if ee, ok := runErr.(*exec.ExitError); ok && ee.ExitCode() == 2 {
			return resp, runErr
		}
		return resp, nil
	}
	if s := strings.TrimSpace(stderr.String()); s != "" {
		// Unexpected stderr with success exit — keep on Error for visibility.
		resp.Error = s
	}
	return resp, nil
}
```
