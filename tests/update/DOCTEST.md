# Update Package Tests

Doc-style tests for `install.HandleUpdate` and `install.HandleUpdateMany`. Update
refreshes **only** directories where `SKILL.md` already exists; missing installs
are skipped with no new files. Batch `HandleUpdateMany` prints
`skill not installed: <cli-name>` per registry entry with no `SKILL.md`; single
`HandleUpdate` stays silent when not installed.

# DSN (Domain Specific Notion)

Participants:

- **Update handler** — parses update-only flags (`--global`, target tool flags,
  `--dry-run`, `--help`), resolves target directories, and never calls
  `HandleInstall`.
- **Target resolver** — maps flags to the same directory layout as install
  (`.agents/skills/<name>`, `.cursor/skills/<name>`, etc., optionally under
  `$HOME` when `--global` is set).
- **Install probe** — treats a target as installed when `<dir>/SKILL.md` exists;
  otherwise the directory is skipped with no log line.
- **InstallTo** — when a target is installed, compares embedded skill content
  (and extra files) to disk; prints `Skill is up to date`, `Update skill at`, or
  dry-run variants; overwrites when content drifted (`noOverride` is always false
  for update).

Behaviors:

- **Single skill** — `HandleUpdate(opts, args)` updates one skill definition
  across all resolved targets that are installed.
- **Registry batch** — `HandleUpdateMany(skills, args)` walks each registry
  entry in stable CLI-name order; when none of its resolved targets have
  `SKILL.md`, prints `skill not installed: <name>` (`UpdateSkill.Name`); when
  any target is installed, runs `InstallTo` only on installed dirs (single-skill
  `HandleUpdate` stays silent when nothing is installed).

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
    │       └── help
    └── many/                         # HandleUpdateMany
        └── registry/
            ├── none-installed-hint
            ├── none-installed-global-hint
            ├── skip-uninstalled
            ├── mixed-reports-both
            └── all-up-to-date
```

## Test Index

| Leaf | Description |
|------|-------------|
| `api/single/not-installed/skip-silent` | No `SKILL.md` → no stdout, no directories created |
| `api/single/installed/up-to-date` | Installed + matching content → `Skill is up to date` |
| `api/single/installed/content-outdated/restores-skill` | Drifted `SKILL.md` → `Update skill at`, content restored |
| `api/single/installed/content-outdated/dry-run` | Drift + `--dry-run` → `[dry-run]` lines, file unchanged |
| `api/single/target-flags/partial-dirs` | Installed under `--codex` only; `--codex --opencode` update touches codex only |
| `api/single/global-scope/global-home` | `--global` install + update; global tree updated, local default untouched |
| `api/single/flags/help` | `--help` prints update usage (location flags, `--dry-run`) |
| `api/many/registry/none-installed-hint` | No installs → per-skill `skill not installed` lines (no scope hint) |
| `api/many/registry/none-installed-global-hint` | `--global`, no installs → not-installed line per registry skill |
| `api/many/registry/skip-uninstalled` | One installed + one missing → up-to-date + `skill not installed: skill-beta` |
| `api/many/registry/mixed-reports-both` | Mixed batch → both install-status line types in CLI-name order |
| `api/many/registry/all-up-to-date` | Two installed skills, current content → two `up to date` lines |

## How to Run

```sh
doctest vet ./tests/update
doctest test -v ./tests/update
```

## Version

0.0.2

```go
import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xhd2015/skills/install"
)

type FileMutate struct {
	RelPath string
	Content string
}

type PreInstall struct {
	Opts install.InstallOptions
	Args []string
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
}

func Run(t *testing.T, req *Request) (*Response, error) {
	tmpDir := t.TempDir()

	prevWD, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	defer os.Chdir(prevWD)

	if err := os.Chdir(tmpDir); err != nil {
		return nil, err
	}

	if req.UseGlobalHome {
		t.Setenv("HOME", t.TempDir())
	}

	for i, m := range req.PostInstallMutate {
		if strings.HasPrefix(m.RelPath, "$HOME/") {
			home, err := os.UserHomeDir()
			if err != nil {
				return nil, err
			}
			req.PostInstallMutate[i].RelPath = filepath.Join(home, strings.TrimPrefix(m.RelPath, "$HOME/"))
		}
	}

	for _, step := range req.PreInstalls {
		if err := install.HandleInstall(step.Opts, step.Args); err != nil {
			return nil, err
		}
	}

	for _, m := range req.PostInstallMutate {
		if err := os.MkdirAll(filepath.Dir(m.RelPath), 0755); err != nil {
			return nil, err
		}
		if err := os.WriteFile(m.RelPath, []byte(m.Content), 0644); err != nil {
			return nil, err
		}
	}

	oldStdout := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		return nil, err
	}

	stdoutCh := make(chan []byte, 1)
	readErrCh := make(chan error, 1)
	go func() {
		data, readErr := io.ReadAll(reader)
		stdoutCh <- data
		readErrCh <- readErr
	}()

	os.Stdout = writer
	var runErr error
	if req.UseMany {
		runErr = install.HandleUpdateMany(req.ManySkills, req.UpdateArgs)
	} else {
		runErr = install.HandleUpdate(req.SingleOpts, req.UpdateArgs)
	}
	os.Stdout = oldStdout
	writer.Close()

	data := <-stdoutCh
	if readErr := <-readErrCh; readErr != nil {
		return nil, readErr
	}
	reader.Close()

	resp := &Response{
		Stdout:  string(data),
		WorkDir: tmpDir,
	}
	if runErr != nil {
		resp.Error = runErr.Error()
	}
	return resp, nil
}
```