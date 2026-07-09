# Install Package Tests

Tests for the `install` package's `HandleInstall` function — verifying output messages,
filesystem side effects, flag interactions, and extra file path validation.

## Decision Tree

```
tests/install/
├── behavior/                          # Install behavior and output messages
│   ├── fresh-install/                 # New directory → "Installed skill to:"
│   ├── overwrite-existing/            # Existing directory → "Update skill at"
│   ├── force-overrides-no-override/   # --force overrides --no-override
│   ├── no-override-empty-dir/         # --no-override on empty dir → installs normally
│   └── cursor-flag/                   # --cursor installs to .cursor/skills/<name>
└── extra-files-validation/            # Extra file path validation errors
    ├── dot-path/                      # Path "." → error
    ├── dotdot-path/                   # Path ".." → error
    └── skill-md-path/                 # Path "SKILL.md" → error
```

## Test Index

| # | Test Leaf | Description |
|---|-----------|-------------|
| 1 | behavior/fresh-install | Fresh install to non-existent directory prints "Installed skill to:" |
| 2 | behavior/overwrite-existing | Overwriting existing directory prints "Update skill at" |
| 3 | behavior/force-overrides-no-override | --force --no-override overwrites without confirmation |
| 4 | behavior/no-override-empty-dir | --no-override on empty dir installs normally (no abort) |
| 5 | behavior/cursor-flag | --cursor installs to .cursor/skills/<name> (local, non-global) |
| 6 | extra-files-validation/dot-path | Extra file path "." returns error |
| 7 | extra-files-validation/dotdot-path | Extra file path ".." returns error |
| 8 | extra-files-validation/skill-md-path | Extra file path "SKILL.md" returns error |

## How to Run

```sh
doctest test -v ./tests/install
```

## Version

0.0.2

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
