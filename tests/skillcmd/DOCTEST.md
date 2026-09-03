# skillcmd Framework Tests

Doc-style tests for package `github.com/xhd2015/skills/skillcmd` — the single home
for skill CLI parse, show/header, install/update, single-skill hosts (flat and
nested tree), and multi-skill registry behavior.

These tests import `skillcmd` directly. They stay **RED** until the package exists
and implements the public surface exercised by `Run`.

# DSN (Domain Specific Notion)

Participants:

- **Caller** — a skill CLI (or this harness) passes argv after the `skill` /
  `skills` token into skillcmd helpers.
- **ParseSkillArgs** — scans args for action flags only (`--show`, `--install`,
  `--list` / `-l`, optional `--header`). Remaining tokens stay in `Rest` (topic
  path, skill name, install flags like `--global`, install dir). Rejects zero or
  multiple actions and rejects combining `--show` with `--install`.
- **SingleSkill** — one skill definition (`Name`, `RootContent`, optional
  `TreeFS` of nested `path/TOPIC.md`, optional `ExtraFiles`). `Handle` runs
  list/show/install for that skill; nested show maps `a/b` → `a/b/TOPIC.md` and
  rejects empty / `.` / `..` segments. Root index remains `SKILL.md` only.
- **Registry** — ordered list of registered skills. `HandleSkill` shows/installs
  by name (both flag orders); `HandleSkills` lists or runs batch update.
- **File header** — `GetHeader` / `ParseHeader` / `FormatHeaderWithDelimiters`
  extract or re-wrap YAML frontmatter (migrated from `skill_file`).
- **Install/Update** — `HandleInstall` / `HandleUpdate` / `HandleUpdateMany`
  write or refresh `.<tool>/skills/<name>/` layouts (migrated from `install`).
  Batch update prints polished `name  not installed` (and a trailing summary)
  when no `SKILL.md` exists. Install inventory treats unplanned on-disk files as
  orphans (delete) and uses incremental create/update/delete (see
  `tests/install` inventory leaves).

Behaviors:

- Actions are **flags only** — no word subcommands `show` / `install`.
- Both orders work: `--show <path>` and `<path> --show` (path lands in Rest).
- Default install target: `.agents/skills/<SkillDirName>/`.
- Explicit `--dir` / positional `<dir>`: smart layout via
  `ResolveExplicitSkillDir` — basename `skills` always nests as
  `<dir>/<name>`; existing `SKILL.md` or matching basename is the skill root;
  otherwise nest. `--dir` is exclusive with `<dir>` and with
  `--cursor/--codex/--opencode/--general-agents`.
- Shape-3 `SingleSkill`: a known topic token before `--install` (show-style
  order) is peeled when `--dir` or another positional destination is present;
  install still copies the whole skill.
- Nested install ExtraFiles use `path/TOPIC.md` (TreeFS load/list/collect), not
  nested `path/SKILL.md` and not `topics/*.md`. Nested `SKILL.md` if present in
  TreeFS is ignored for topic discovery.

## Decision Tree

```text
skillcmd/
├── parse/                         # ParseSkillArgs (argv classification)
│   ├── show-flag-only
│   ├── show-path-before-flag
│   ├── show-flag-before-path
│   ├── show-header
│   ├── install-flag
│   ├── list-flag
│   ├── list-short-flag
│   ├── both-show-and-install
│   └── missing-action
├── single-flat/                   # SingleSkill without TreeFS
│   ├── show-root
│   ├── show-header
│   ├── list
│   └── install-dry-run
├── single-tree/                   # SingleSkill with nested path/TOPIC.md
│   ├── show-nested-path
│   ├── reject-dotdot
│   ├── install-extra-files
│   ├── install-topic-before/      # <topic> --install peels topic
│   │   ├── dir-flag
│   │   └── positional
│   └── list-topics
├── multi/                      # multi-skill host
│   ├── list-skills
│   ├── show-by-name/
│   │   ├── flag-before-name
│   │   └── name-before-flag
│   ├── install-named
│   └── update-many/
│       ├── skip-missing
│       └── updates-installed
├── file-header/                   # skillcmd file/header APIs
│   ├── get-header
│   └── format-header
└── install-compat/                # HandleInstall via skillcmd
    ├── fresh-default
    ├── positional-collection      # <dir> basename skills → nest
    └── dir-flag/
        ├── collection             # --dir …/skills → nest
        ├── matching-basename      # --dir …/<name> → direct
        ├── conflict-positional    # --dir + <dir> → error
        └── conflict-cursor        # --dir + --cursor → error
```

## Test Index

| Leaf | Description |
|------|-------------|
| `parse/show-flag-only` | `--show` → action show, empty rest |
| `parse/show-path-before-flag` | `flags-parsing/types --show` → show + path in rest |
| `parse/show-flag-before-path` | `--show flags-parsing/types` → show + path in rest |
| `parse/show-header` | `--show --header` → show with Header true |
| `parse/install-flag` | `--install --global` → install, rest keeps `--global` |
| `parse/list-flag` | `--list` → action list |
| `parse/list-short-flag` | `-l` → action list |
| `parse/both-show-and-install` | `--show --install` → error |
| `parse/missing-action` | bare `foo` → error (no action flag) |
| `single-flat/show-root` | `--show` prints RootContent |
| `single-flat/show-header` | `--show --header` prints delimiters + name, no body |
| `single-flat/list` | `--list` prints skill Name |
| `single-flat/install-dry-run` | `--install --dry-run` mentions `.agents/skills/<name>` |
| `single-tree/show-nested-path` | `--show a/b` prints nested TOPIC.md body |
| `single-tree/reject-dotdot` | `--show ../x` errors on invalid segment |
| `single-tree/install-extra-files` | install writes `skill-cli/TOPIC.md` (not nested SKILL.md / topics/*) |
| `single-tree/install-topic-before/dir-flag` | `skill-cli --install --dir vendor/skills` peels topic |
| `single-tree/install-topic-before/positional` | `skill-cli --install vendor/skills` peels topic |
| `single-tree/list-topics` | `--list` lists skill name + topics from `**/TOPIC.md` |
| `multi/list-skills` | `--list` lists registered names (+ description) |
| `multi/show-by-name/flag-before-name` | `--show foo` prints foo content |
| `multi/show-by-name/name-before-flag` | `foo --show` prints same content |
| `multi/install-named` | `--install foo --dry-run` targets foo skill dir |
| `multi/update-many/skip-missing` | not installed → polished `name  not installed` + summary |
| `multi/update-many/updates-installed` | existing SKILL.md refreshed; polished `updated` status |
| `file-header/get-header` | GetHeader returns inner YAML without delimiters |
| `file-header/format-header` | FormatHeaderWithDelimiters wraps with `---` |
| `install-compat/fresh-default` | HandleInstall default dir succeeds under skillcmd |
| `install-compat/positional-collection` | positional `vendor/skills` nests to `…/demo-skill` |
| `install-compat/dir-flag/collection` | `--dir vendor/skills` nests; no collection-root SKILL.md |
| `install-compat/dir-flag/matching-basename` | `--dir out/demo-skill` installs directly |
| `install-compat/dir-flag/conflict-positional` | `--dir` + positional → error |
| `install-compat/dir-flag/conflict-cursor` | `--dir` + `--cursor` → error |

## How to Run

```sh
doctest vet ./tests/skillcmd
doctest test -v ./tests/skillcmd
```

## Version

0.0.2

```go
import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"testing/fstest"

	"github.com/xhd2015/doctest/session"
	"github.com/xhd2015/skills/skillcmd"
)

// processMu serializes process-global mutations in Run (chdir, os.Stdout,
// HOME). Doctest leaves use t.Parallel(); skillcmd product APIs resolve relative
// skill paths from the process cwd and print via os.Stdout.
var processMu sync.Mutex


// Mode selects which skillcmd surface Run exercises.
type Mode string

const (
	ModeParse         Mode = "parse"
	ModeSingle        Mode = "single"
	ModeRegistry      Mode = "registry"
	ModeFileHeader    Mode = "file-header"
	ModeInstallCompat Mode = "install-compat"
)

// FileOp selects which header helper to call in ModeFileHeader.
type FileOp string

const (
	FileOpGetHeader    FileOp = "get-header"
	FileOpFormatHeader FileOp = "format-header"
)

// RegistryCmd selects Registry entry point.
type RegistryCmd string

const (
	RegistryCmdSkill  RegistryCmd = "skill"
	RegistryCmdSkills RegistryCmd = "skills"
)

type Request struct {
	Mode Mode
	Args []string

	// SingleSkill / shared skill identity
	SkillName   string
	RootContent string
	// TreeFiles maps slash paths (e.g. "a/b/TOPIC.md") → file content.
	// When non-empty, SingleSkill.TreeFS is built from this map.
	TreeFiles map[string]string
	// ExtraFiles overrides auto-collection from TreeFiles when non-nil.
	ExtraFiles []skillcmd.InstallFile
	Usage      string

	// Registry
	RegistrySkills []skillcmd.RegisteredSkill
	RegistryCmd    RegistryCmd

	// File header
	Content string
	FileOp  FileOp

	// Install-compat / update helpers
	InstallOpts       skillcmd.InstallOptions
	PreInstall        bool
	PreInstallOpts    skillcmd.InstallOptions
	PreInstallArgs    []string
	PostInstallMutate map[string]string // rel path → content after pre-install

	UseWorkDir    bool
	UseGlobalHome bool
}

type Response struct {
	// Parse result
	Action string
	Header bool
	Rest   []string

	// Captured handler output / errors
	Stdout string
	Error  string

	// File-header
	HeaderText string
	HeaderErr  string
	Formatted  string
	FormatErr  string

	WorkDir string
	HomeDir string
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	t.Helper()
	_ = d
	if req.Mode == "" {
		return nil, fmt.Errorf("req.Mode is required")
	}

	// Serialize process-global side effects used by install/update/show handlers.
	processMu.Lock()
	defer processMu.Unlock()

	resp := &Response{}

	if req.UseWorkDir {
		workDir := t.TempDir()
		resp.WorkDir = workDir
		prevWD, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		if err := os.Chdir(workDir); err != nil {
			return nil, err
		}
		defer func() { _ = os.Chdir(prevWD) }()
	}

	if req.UseGlobalHome {
		homeDir := t.TempDir()
		resp.HomeDir = homeDir
		t.Setenv("HOME", homeDir)
	}

	switch req.Mode {
	case ModeParse:
		parsed, err := skillcmd.ParseSkillArgs(req.Args)
		if err != nil {
			resp.Error = err.Error()
			return resp, nil
		}
		resp.Action = string(parsed.Action)
		resp.Header = parsed.Header
		resp.Rest = append([]string(nil), parsed.Rest...)
		return resp, nil

	case ModeSingle:
		stdout, errStr := captureStdout(t, func() error {
			sk := buildSingleSkill(req)
			return sk.Handle(req.Args)
		})
		resp.Stdout = stdout
		resp.Error = errStr
		return resp, nil

	case ModeRegistry:
		if req.PreInstall {
			if err := skillcmd.HandleInstall(req.PreInstallOpts, req.PreInstallArgs); err != nil {
				return nil, fmt.Errorf("pre-install: %w", err)
			}
			for rel, content := range req.PostInstallMutate {
				p := filepath.Join(resp.WorkDir, filepath.FromSlash(rel))
				if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
					return nil, err
				}
				if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
					return nil, err
				}
			}
		}
		reg := &skillcmd.Registry{Skills: req.RegistrySkills}
		stdout, errStr := captureStdout(t, func() error {
			switch req.RegistryCmd {
			case RegistryCmdSkills:
				return reg.HandleSkills(req.Args)
			default:
				return reg.HandleSkill(req.Args)
			}
		})
		resp.Stdout = stdout
		resp.Error = errStr
		return resp, nil

	case ModeFileHeader:
		switch req.FileOp {
		case FileOpGetHeader:
			h, err := skillcmd.GetHeader(req.Content)
			resp.HeaderText = h
			if err != nil {
				resp.HeaderErr = err.Error()
			}
		case FileOpFormatHeader:
			out, err := skillcmd.FormatHeaderWithDelimiters(req.Content)
			resp.Formatted = out
			if err != nil {
				resp.FormatErr = err.Error()
			}
		default:
			return nil, fmt.Errorf("unknown FileOp %q", req.FileOp)
		}
		return resp, nil

	case ModeInstallCompat:
		stdout, errStr := captureStdout(t, func() error {
			return skillcmd.HandleInstall(req.InstallOpts, req.Args)
		})
		resp.Stdout = stdout
		resp.Error = errStr
		return resp, nil

	default:
		return nil, fmt.Errorf("unknown mode %q", req.Mode)
	}
}

func buildSingleSkill(req *Request) *skillcmd.SingleSkill {
	sk := &skillcmd.SingleSkill{
		Name:        req.SkillName,
		RootContent: req.RootContent,
		Usage:       req.Usage,
	}
	// Explicit ExtraFiles override only when the leaf sets the field.
	// When nil, SingleSkill.install derives from TreeFS via collectTreeSkillFiles
	// (must discover nested TOPIC.md, not nested SKILL.md).
	if req.ExtraFiles != nil {
		sk.ExtraFiles = req.ExtraFiles
	}
	if len(req.TreeFiles) > 0 {
		m := fstest.MapFS{}
		for p, content := range req.TreeFiles {
			m[p] = &fstest.MapFile{Data: []byte(content)}
		}
		sk.TreeFS = fs.FS(m)
	}
	return sk
}

func captureStdout(t *testing.T, fn func() error) (string, string) {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	old := os.Stdout
	os.Stdout = w
	runErr := fn()
	_ = w.Close()
	os.Stdout = old
	data, readErr := io.ReadAll(r)
	_ = r.Close()
	if readErr != nil {
		t.Fatalf("read stdout: %v", readErr)
	}
	errStr := ""
	if runErr != nil {
		errStr = runErr.Error()
	}
	return string(data), errStr
}
```
