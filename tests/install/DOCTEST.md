# Install Package Tests

Tests for the `install` package's `HandleInstall` function — verifying output messages,
filesystem side effects, flag interactions, extra file path validation, and
skill-dir inventory sync (rsync-like create/update/delete).

# DSN (Domain Specific Notion)

Participants:

- **Caller** — test harness (or a skill CLI) invokes `install.HandleInstall` with
  skill content, optional `ExtraFiles`, and install flags / target dir.
- **HandleInstall** — parses flags (`--dry-run`, `--no-override`, `--force`,
  target tool flags, `--global`), resolves one or more skill target directories,
  and runs inventory sync per target.
- **Install plan** — desired regular-file set =
  `{ "SKILL.md": SkillContent }` ∪ ExtraFiles (relative paths under the skill dir).
  Nested topics use paths like `a/TOPIC.md` (not nested `SKILL.md`).
- **Skill directory** — on-disk tree at the resolved target (e.g. `example-skill/`
  or `.agents/skills/<name>/`). Only files under this root are inventory-managed.
- **Inventory sync** — compares the plan to **all regular files** under the skill
  dir (recursive). Missing/different planned paths → create or update; unplanned
  on-disk files → delete; empty dirs left after deletes are removed best-effort.
  Does **not** `RemoveAll` the whole skill dir on normal updates.

Behaviors:

- **Up to date** — all planned paths match by content **and** no unplanned files
  on disk → print `Skill is up to date: <absDir>\n` and write nothing.
- **Work needed** — print header then per-file actions (absolute paths preferred):
  - Header: `Installed skill to: <absDir>` only if the skill dir did **not** exist
    before this install; otherwise `Update skill at <absDir>`.
  - Detail lines, stable order: all `create: ` (sorted), then `update: `, then
    `delete: ` (prefix is `create: `/`update: `/`delete: ` with one space after colon).
  - No MD5 hashes in stdout; last line ends with `\n`.
- **Dry-run** — same messages with `[dry-run] ` prefix; no write/delete.
  Up to date dry-run: `[dry-run] Skill is up to date: <absDir>\n` only.
- **`--no-override`** — if skill dir is non-empty and any create/update/delete is
  required, confirm (TTY) or abort (non-TTY). Empty or missing dir: no confirmation.
- **Extra path validation** — reject `.`, `..`, absolute paths, and extra path
  `SKILL.md` before any inventory work.

## Decision Tree

```
tests/install/
├── behavior/                          # Install behavior and output messages
│   ├── fresh-install/                 # New directory → "Installed skill to:"
│   ├── overwrite-existing/            # Existing directory → "Update skill at"
│   ├── force-overrides-no-override/   # --force overrides --no-override
│   ├── no-override-empty-dir/         # --no-override on empty dir → no abort
│   └── cursor-flag/                   # --cursor installs to .cursor/skills/<name>
├── extra-files-validation/            # Extra file path validation errors
│   ├── dot-path/                      # Path "." → error
│   ├── dotdot-path/                   # Path ".." → error
│   └── skill-md-path/                 # Path "SKILL.md" → error
└── inventory/                         # Skill-dir inventory sync (rsync-like)
    ├── clean-match/                   # disk == plan, no orphans
    │   ├── apply-up-to-date/          # apply → Skill is up to date; no writes
    │   └── dry-run-up-to-date/        # --dry-run → [dry-run] Skill is up to date
    ├── orphan/                        # plan matches + unplanned on-disk files
    │   ├── apply-deletes/             # delete: orphan; file removed
    │   └── dry-run-preserves/         # dry-run delete line; file still on disk
    ├── plan-changed/                  # planned set or content differs
    │   ├── content-update/            # update: existing path; new content
    │   ├── drop-nested/               # delete: nested path dropped from plan
    │   └── rename-skill-to-topic/     # delete: a/SKILL.md + create: a/TOPIC.md
    └── fresh/                         # skill dir missing
        └── create-with-extras/        # Installed + create: for root and extras
```

## Test Index

| # | Test Leaf | Description |
|---|-----------|-------------|
| 1 | behavior/fresh-install | Fresh install to non-existent directory prints "Installed skill to:" |
| 2 | behavior/overwrite-existing | Overwriting existing directory prints "Update skill at" |
| 3 | behavior/force-overrides-no-override | --force --no-override overwrites without confirmation |
| 4 | behavior/no-override-empty-dir | --no-override on empty dir installs without abort (dir existed → Update header) |
| 5 | behavior/cursor-flag | --cursor installs to .cursor/skills/<name> (local, non-global) |
| 6 | extra-files-validation/dot-path | Extra file path "." returns error |
| 7 | extra-files-validation/dotdot-path | Extra file path ".." returns error |
| 8 | extra-files-validation/skill-md-path | Extra file path "SKILL.md" returns error |
| 9 | inventory/clean-match/apply-up-to-date | Clean disk==plan → up to date; no file changes |
| 10 | inventory/clean-match/dry-run-up-to-date | Dry-run clean match → `[dry-run] Skill is up to date` only |
| 11 | inventory/orphan/apply-deletes | Orphan leftover → not up to date; `delete:`; file gone |
| 12 | inventory/orphan/dry-run-preserves | Dry-run orphan → delete line; file remains |
| 13 | inventory/plan-changed/content-update | Content differs → `update:`; new content on disk |
| 14 | inventory/plan-changed/drop-nested | Nested file removed from plan → `delete:`; gone |
| 15 | inventory/plan-changed/rename-skill-to-topic | a/SKILL.md → a/TOPIC.md as delete+create |
| 16 | inventory/fresh/create-with-extras | Missing dir → `Installed skill to:` + sorted `create:` lines |

## How to Run

```sh
doctest vet ./tests/install
doctest test -v ./tests/install
```

## Version

0.0.2

```go
import (
	"io"
	"os"
	"path/filepath"
	"syscall"
	"testing"

	"github.com/xhd2015/doctest/session"
	"github.com/xhd2015/skills/install"
)

// withProcessLock serializes process-global mutations across all doctest trees
// in the same go test process (shared flock path).
func withProcessLock(fn func() error) error {
	lockPath := filepath.Join(os.TempDir(), "skills-doctest-process.lock")
	f, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		return err
	}
	defer syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
	return fn()
}


type Request struct {
	SkillDirName string
	SkillContent string
	Args         []string
	ExtraFiles   []install.InstallFile

	// PreExistingDir is a directory path (relative to working dir) to create before calling HandleInstall.
	PreExistingDir   string
	PreExistingFiles []PreExistingFile

	// NonInteractive replaces stdin with an empty file so confirmOverwrite returns false.
	NonInteractive bool

	// UseGlobalHome sets HOME to a temp dir, used for --global flag tests.
	UseGlobalHome bool
}

type PreExistingFile struct {
	// Name is relative to PreExistingDir; may include nested path segments
	// (slash- or OS-separated). Parents are created with MkdirAll.
	Name    string
	Content string
}

type Response struct {
	Stdout  string
	Error   string
	WorkDir string // absolute path to the working temp directory
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	_ = d
	var resp *Response
	err := withProcessLock(func() error {
		var runErr error
		resp, runErr = runLocked(t, req)
		return runErr
	})
	return resp, err
}

func runLocked(t *testing.T, req *Request) (*Response, error) {
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
		homeDir := t.TempDir()
		prevHome, hadHome := os.LookupEnv("HOME")
		if err := os.Setenv("HOME", homeDir); err != nil {
			return nil, err
		}
		defer func() {
			if hadHome {
				_ = os.Setenv("HOME", prevHome)
			} else {
				_ = os.Unsetenv("HOME")
			}
		}()
	}

	// Set up pre-existing directory and files (nested Names get parent MkdirAll).
	if req.PreExistingDir != "" {
		if err := os.MkdirAll(req.PreExistingDir, 0755); err != nil {
			return nil, err
		}
		for _, f := range req.PreExistingFiles {
			dest := filepath.Join(req.PreExistingDir, filepath.FromSlash(f.Name))
			if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
				return nil, err
			}
			if err := os.WriteFile(dest, []byte(f.Content), 0644); err != nil {
				return nil, err
			}
		}
	}

	// Replace stdin for non-interactive mode
	if req.NonInteractive {
		stdinPath := filepath.Join(t.TempDir(), "stdin")
		if err := os.WriteFile(stdinPath, nil, 0644); err != nil {
			return nil, err
		}
		stdin, err := os.Open(stdinPath)
		if err != nil {
			return nil, err
		}
		defer stdin.Close()
		prevStdin := os.Stdin
		os.Stdin = stdin
		defer func() { os.Stdin = prevStdin }()
	}

	// Capture stdout
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
	runErr := install.HandleInstall(install.InstallOptions{
		SkillDirName: req.SkillDirName,
		SkillContent: req.SkillContent,
		ExtraFiles:   req.ExtraFiles,
	}, req.Args)
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
