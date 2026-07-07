package main

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/xhd2015/skills/playwrightdebug"
)

func TestHelpText(t *testing.T) {
	if !strings.Contains(help, "playwright-debug") {
		t.Error("help text missing playwright-debug reference")
	}
	if !strings.Contains(help, "run") {
		t.Error("help text missing run command")
	}
	if !strings.Contains(help, "skill install") {
		t.Error("help text missing skill install command")
	}
	if !strings.Contains(help, "skill show") {
		t.Error("help text missing skill show command")
	}
}

func TestCacheDir(t *testing.T) {
	dir := playwrightdebug.DefaultCacheDir()
	if !strings.HasSuffix(dir, ".playwright-debug/node_package") {
		t.Errorf("unexpected cache dir: %s", dir)
	}
}

func TestHandleNoArgs(t *testing.T) {
	output, err := captureStdout(t, func() error {
		return handle(nil)
	})
	if err != nil {
		t.Fatalf("handle(nil): %v", err)
	}
	if !strings.Contains(output, "Usage:") {
		t.Errorf("no-args output missing Usage: %s", output)
	}
}

func TestHandleHelpShort(t *testing.T) {
	output, err := captureStdout(t, func() error {
		return handle([]string{"-h"})
	})
	if err != nil {
		t.Fatalf("handle(-h): %v", err)
	}
	if !strings.Contains(output, "Usage:") {
		t.Errorf("help output missing Usage: %s", output)
	}
}

func TestHandleHelpLong(t *testing.T) {
	output, err := captureStdout(t, func() error {
		return handle([]string{"--help"})
	})
	if err != nil {
		t.Fatalf("handle(--help): %v", err)
	}
	if !strings.Contains(output, "Usage:") {
		t.Errorf("help output missing Usage: %s", output)
	}
}

func TestHandleRunMissingScript(t *testing.T) {
	err := handle([]string{"run"})
	if err == nil {
		t.Fatal("expected error for run without script")
	}
	msg := strings.ToLower(err.Error())
	if !strings.Contains(msg, "require") {
		t.Errorf("unexpected error: %v", err)
	}
	if !strings.Contains(msg, "file") && !strings.Contains(msg, ".js") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestHandleRunHelpInArgs(t *testing.T) {
	output, err := captureStdout(t, func() error {
		return handle([]string{"--help", "console.log(1)"})
	})
	if err != nil {
		t.Fatalf("handle(--help with args): %v", err)
	}
	if strings.Contains(output, "Usage:") {
		t.Errorf("CLI help must not be shown when --help appears with other args, got: %s", output)
	}
}

func TestHandleSkillShow(t *testing.T) {
	output, err := captureStdout(t, func() error {
		return handle([]string{"skill", "show"})
	})
	if err != nil {
		t.Fatalf("handle(skill show): %v", err)
	}
	if !strings.Contains(output, "playwright-debug") {
		t.Errorf("skill show output missing skill name: %s", output)
	}
}

func TestHandleSkillUnknown(t *testing.T) {
	err := handle([]string{"skill", "unknown"})
	if err == nil {
		t.Fatal("expected error for unknown skill sub-command")
	}
	if !strings.Contains(err.Error(), "unknown skill sub-command") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestHandleSkillNoSubcommand(t *testing.T) {
	err := handle([]string{"skill"})
	if err == nil {
		t.Fatal("expected error for skill without sub-command")
	}
	if !strings.Contains(err.Error(), "unknown skill sub-command") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestHandleSkillInstallDryRun(t *testing.T) {
	output, err := captureStdout(t, func() error {
		return handle([]string{"skill", "install", "--dry-run"})
	})
	if err != nil {
		t.Fatalf("handle(skill install --dry-run): %v", err)
	}
	if !strings.Contains(output, "[dry-run]") {
		t.Errorf("expected dry-run output, got: %s", output)
	}
}

func TestHandleSkillInstallDirDryRun(t *testing.T) {
	output, err := captureStdout(t, func() error {
		return handle([]string{"skill", "install", "/tmp/test-playwright-debug", "--dry-run"})
	})
	if err != nil {
		t.Fatalf("handle(skill install /tmp/... --dry-run): %v", err)
	}
	if !strings.Contains(output, "[dry-run]") {
		t.Errorf("expected dry-run output, got: %s", output)
	}
	if !strings.Contains(output, "/tmp/test-playwright-debug") {
		t.Errorf("expected output to mention target dir, got: %s", output)
	}
}

func TestHandleSkillInstallCursorDryRun(t *testing.T) {
	output, err := captureStdout(t, func() error {
		return handle([]string{"skill", "install", "--cursor", "--dry-run"})
	})
	if err != nil {
		t.Fatalf("handle(skill install --cursor --dry-run): %v", err)
	}
	if !strings.Contains(output, "[dry-run]") {
		t.Errorf("expected dry-run output, got: %s", output)
	}
}

func captureStdout(t *testing.T, fn func() error) (string, error) {
	t.Helper()
	oldStdout := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("create stdout pipe: %v", err)
	}

	stdoutCh := make(chan []byte, 1)
	readErrCh := make(chan error, 1)
	go func() {
		data, readErr := io.ReadAll(reader)
		stdoutCh <- data
		readErrCh <- readErr
	}()

	os.Stdout = writer
	runErr := fn()
	os.Stdout = oldStdout
	if err := writer.Close(); err != nil {
		t.Fatalf("close stdout writer: %v", err)
	}
	data := <-stdoutCh
	if err := <-readErrCh; err != nil {
		t.Fatalf("read stdout: %v", err)
	}
	if err := reader.Close(); err != nil {
		t.Fatalf("close stdout reader: %v", err)
	}
	return string(data), runErr
}
