## Preconditions
- The install package's `HandleInstall` function is available at the module root.
- Tests run in an isolated temporary directory.
- Working directory is changed to a fresh temp dir before each test.

## Steps
1. Create a temporary directory and change into it.
2. If `req.UseGlobalHome` is true, set `HOME` to a separate temporary directory.
3. If `req.PreExistingDir` is non-empty, create the directory with the specified pre-existing files.
4. If `req.NonInteractive` is true, replace stdin with an empty file so `confirmOverwrite` returns false without waiting for user input.
5. Capture stdout via an os.Pipe.
6. Call `HandleInstall(InstallOptions{SkillDirName, SkillContent, ExtraFiles, Usage: ""}, req.Args)`.
7. Collect captured stdout and any error returned by HandleInstall into the Response.

## Context
- `HandleInstall` is the entry point that parses flags from `args` and calls `installTo` for each target directory.
- Supported flags: `--cursor`, `--codex`, `--opencode`, `--general-agents`, `--global`, `--no-override`, `--force`, `--dry-run`.
- Output messages:
  - Fresh install (directory does not exist): `"Installed skill to: <absDir>"`
  - Overwrite (directory exists and is replaced): `"Update skill at <absDir>"`
  - Up-to-date (SKILL.md unchanged): `"Skill is up to date: <absDir>"`
  - Dry-run variants: prefixed with `"[dry-run] "`
  - Aborted (user declined confirmation): `"Aborted."`
- Extra file paths are validated in `resolveInstallFiles`: ".", "..", absolute paths, and "SKILL.md" are rejected with `"invalid install file path"` or `"extra install file cannot replace SKILL.md"`.

```go
import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/xhd2015/skills/install"
)

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
	Name    string
	Content string
}

type Response struct {
	Stdout  string
	Error   string
	WorkDir string // absolute path to the working temp directory
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
		homeDir := t.TempDir()
		t.Setenv("HOME", homeDir)
	}

	// Set up pre-existing directory and files
	if req.PreExistingDir != "" {
		if err := os.MkdirAll(req.PreExistingDir, 0755); err != nil {
			return nil, err
		}
		for _, f := range req.PreExistingFiles {
			if err := os.WriteFile(filepath.Join(req.PreExistingDir, f.Name), []byte(f.Content), 0644); err != nil {
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
