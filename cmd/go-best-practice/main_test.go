package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestListTopics(t *testing.T) {
	topics, err := listTopics()
	if err != nil {
		t.Fatalf("listTopics: %v", err)
	}
	topicSet := make(map[string]bool, len(topics))
	for _, tp := range topics {
		topicSet[tp] = true
	}
	expectedTopics := []string{
		"cmd-exec",
		"flags-parsing",
		"flags-parsing/subcommand",
		"flags-parsing/types",
		"kool-create",
		"skill-cli",
	}
	for _, expected := range expectedTopics {
		if !topicSet[expected] {
			t.Errorf("missing topic: %q", expected)
		}
	}
}

func TestReadTopicExistingTopLevel(t *testing.T) {
	content, ok, err := readTopic("kool-create")
	if err != nil {
		t.Fatalf("readTopic(kool-create): %v", err)
	}
	if !ok {
		t.Fatal("expected ok for kool-create")
	}
	if !strings.Contains(content, "kool create") {
		t.Errorf("unexpected content for kool-create: %s", content)
	}
}

func TestReadTopicSkillCLI(t *testing.T) {
	content, ok, err := readTopic("skill-cli")
	if err != nil {
		t.Fatalf("readTopic(skill-cli): %v", err)
	}
	if !ok {
		t.Fatal("expected ok for skill-cli")
	}
	if !strings.Contains(content, "skill install") {
		t.Errorf("unexpected content for skill-cli: %s", content)
	}
}

func TestHandleSkillCLITopic(t *testing.T) {
	output, err := captureStdout(t, func() error {
		return handle([]string{"skill-cli"})
	})
	if err != nil {
		t.Fatalf("handle(skill-cli): %v", err)
	}
	if !strings.Contains(output, "skill install") {
		t.Errorf("skill-cli output missing expected content: %s", output)
	}
}

func TestReadTopicExistingSubTopic(t *testing.T) {
	content, ok, err := readTopic("flags-parsing/types")
	if err != nil {
		t.Fatalf("readTopic(flags-parsing/types): %v", err)
	}
	if !ok {
		t.Fatal("expected ok for flags-parsing/types")
	}
	if !strings.Contains(content, "**T") {
		t.Errorf("unexpected content for flags-parsing/types: %s", content)
	}
}

func TestReadTopicNonExistent(t *testing.T) {
	_, ok, err := readTopic("nonexistent")
	if err != nil {
		t.Fatalf("readTopic(nonexistent): %v", err)
	}
	if ok {
		t.Error("expected not ok for nonexistent topic")
	}
}

func TestReadTopicEmpty(t *testing.T) {
	_, ok, err := readTopic("")
	if err != nil {
		t.Fatalf("readTopic(''): %v", err)
	}
	if ok {
		t.Error("expected not ok for empty topic")
	}
}

func TestValidateSegmentsOK(t *testing.T) {
	tests := []struct {
		segments []string
	}{
		{[]string{"kool-create"}},
		{[]string{"flags-parsing", "types"}},
		{[]string{"a"}},
	}
	for _, tt := range tests {
		if err := validateSegments(tt.segments); err != nil {
			t.Errorf("validateSegments(%v): unexpected error: %v", tt.segments, err)
		}
	}
}

func TestValidateSegmentsReject(t *testing.T) {
	tests := []struct {
		segments []string
	}{
		{[]string{""}},
		{[]string{"."}},
		{[]string{".."}},
		{[]string{"a", ""}},
		{[]string{"a", ".."}},
	}
	for _, tt := range tests {
		if err := validateSegments(tt.segments); err == nil {
			t.Errorf("validateSegments(%v): expected error, got nil", tt.segments)
		}
	}
}

func TestPrintTopicIndex(t *testing.T) {
	output, err := captureStdout(t, func() error {
		return printTopicIndex()
	})
	if err != nil {
		t.Fatalf("printTopicIndex: %v", err)
	}
	if !strings.Contains(output, "Available topics:") {
		t.Errorf("missing header in output: %s", output)
	}
	if !strings.Contains(output, "kool-create") {
		t.Errorf("topics output missing kool-create: %s", output)
	}
	if !strings.Contains(output, "skill-cli") {
		t.Errorf("topics output missing skill-cli: %s", output)
	}
	if !strings.Contains(output, "skill-cli") {
		t.Errorf("missing skill-cli in output: %s", output)
	}
}

func TestHelpText(t *testing.T) {
	if !strings.Contains(help, "go-best-practice") {
		t.Error("help text missing go-best-practice reference")
	}
	if !strings.Contains(help, "install") {
		t.Error("help text missing install command")
	}
	if !strings.Contains(help, "topics") {
		t.Error("help text missing topics command")
	}
}

func TestCollectTopicFiles(t *testing.T) {
	files, err := collectTopicFiles()
	if err != nil {
		t.Fatalf("collectTopicFiles: %v", err)
	}
	if len(files) == 0 {
		t.Fatal("expected at least one topic file")
	}
	hasCmdExec := false
	for _, f := range files {
		if f.Path == "SKILL.md" {
			t.Error("topic files should not include SKILL.md")
		}
		if f.Path == "topics/cmd-exec.md" {
			hasCmdExec = true
		}
		if len(f.Content) == 0 {
			t.Errorf("empty content for %s", f.Path)
		}
	}
	if !hasCmdExec {
		t.Error("missing topics/cmd-exec.md in collected files")
	}
}

func TestHandleHelp(t *testing.T) {
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

func TestHandleTopics(t *testing.T) {
	output, err := captureStdout(t, func() error {
		return handle([]string{"topics"})
	})
	if err != nil {
		t.Fatalf("handle(topics): %v", err)
	}
	if !strings.Contains(output, "Available topics:") {
		t.Errorf("topics output missing header: %s", output)
	}
}

func TestHandleKnownTopic(t *testing.T) {
	output, err := captureStdout(t, func() error {
		return handle([]string{"kool-create"})
	})
	if err != nil {
		t.Fatalf("handle(kool-create): %v", err)
	}
	if !strings.Contains(output, "kool create") {
		t.Errorf("topic output missing expected content: %s", output)
	}
}

func TestHandleUnknownCommand(t *testing.T) {
	err := handle([]string{"unknown-cmd"})
	if err == nil {
		t.Fatal("expected error for unknown command")
	}
	if !strings.Contains(err.Error(), "unknown command or topic") {
		t.Errorf("unexpected error message: %v", err)
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
	if !strings.Contains(output, "Available topics:") {
		t.Errorf("no-args output missing topic index: %s", output)
	}
}

func TestHandleSkillShow(t *testing.T) {
	output, err := captureStdout(t, func() error {
		return handle([]string{"skill", "show"})
	})
	if err != nil {
		t.Fatalf("handle(skill show): %v", err)
	}
	if !strings.Contains(output, "go-best-practice") {
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

func TestHandleVet_NoViolations(t *testing.T) {
	dir := t.TempDir()
	mustWriteFile(t, dir, "main.go", `package main
import "fmt"
func main() { fmt.Println("hello") }
`)
	output, err := captureStdout(t, func() error {
		return handle([]string{"vet", dir})
	})
	if err != nil {
		t.Fatalf("handle(vet): %v", err)
	}
	if strings.TrimSpace(output) != "" {
		t.Errorf("expected no output for clean dir, got: %s", output)
	}
}

func TestHandleVet_StdFlag(t *testing.T) {
	dir := t.TempDir()
	mustWriteFile(t, dir, "main.go", `package main
import "flag"
func main() { flag.Parse() }
`)
	output, err := captureStdout(t, func() error {
		return handle([]string{"vet", dir})
	})
	if err != nil {
		t.Fatalf("handle(vet): %v", err)
	}
	if !strings.Contains(output, "[std-flag]") {
		t.Errorf("expected std-flag violation, got: %s", output)
	}
	if !strings.Contains(output, "flags-parsing") {
		t.Errorf("expected hint about flags-parsing, got: %s", output)
	}
}

func TestHandleVet_FileLength(t *testing.T) {
	dir := t.TempDir()
	content := "package main\n" + strings.Repeat("// line\n", 600)
	mustWriteFile(t, dir, "main.go", content)

	output, err := captureStdout(t, func() error {
		return handle([]string{"vet", "--file-max-lines", "500", dir})
	})
	if err != nil {
		t.Fatalf("handle(vet): %v", err)
	}
	if !strings.Contains(output, "[file-length]") {
		t.Errorf("expected file-length violation, got: %s", output)
	}
	if !strings.Contains(output, "exceeds maximum of 500") {
		t.Errorf("expected threshold message, got: %s", output)
	}
}

func TestHandleVet_JsonOutput(t *testing.T) {
	dir := t.TempDir()
	mustWriteFile(t, dir, "main.go", `package main
import "flag"
func main() { flag.Parse() }
`)
	output, err := captureStdout(t, func() error {
		return handle([]string{"vet", "--json", dir})
	})
	if err != nil {
		t.Fatalf("handle(vet --json): %v", err)
	}
	if !strings.Contains(output, "\"Checker\"") {
		t.Errorf("expected JSON output, got: %s", output)
	}
}

func TestHandleVet_Exclude(t *testing.T) {
	dir := t.TempDir()
	mustWriteFile(t, dir, "main.go", `package main
import "flag"
func main() { flag.Parse() }
`)
	output, err := captureStdout(t, func() error {
		return handle([]string{"vet", "--exclude", "std-flag", dir})
	})
	if err != nil {
		t.Fatalf("handle(vet --exclude): %v", err)
	}
	if strings.Contains(output, "[std-flag]") {
		t.Errorf("expected std-flag to be excluded, got: %s", output)
	}
}

func TestHandleVet_NoArgs(t *testing.T) {
	output, err := captureStdout(t, func() error {
		return handle([]string{"vet"})
	})
	if err != nil {
		t.Fatalf("handle(vet no-args): %v", err)
	}
	if strings.TrimSpace(output) != "" {
		t.Errorf("expected no output for vet on current dir (no violations), got: %s", output)
	}
}

func TestHandleVet_UnknownFlag(t *testing.T) {
	err := handle([]string{"vet", "--unknown"})
	if err == nil {
		t.Fatal("expected error for unknown flag")
	}
	if !strings.Contains(err.Error(), "unknown flag") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestHandleVet_BadFileMaxLines(t *testing.T) {
	tests := []struct {
		flag string
	}{
		{"--file-max-lines=abc"},
		{"--file-max-lines=-1"},
	}
	for _, tt := range tests {
		t.Run(tt.flag, func(t *testing.T) {
			err := handle([]string{"vet", tt.flag, "."})
			if err == nil {
				t.Fatal("expected error for bad --file-max-lines")
			}
		})
	}
}

func TestHandleVet_DotDotDot(t *testing.T) {
	root := t.TempDir()
	mustWriteFile(t, root, "go.mod", "module testmod\ngo 1.24\n")
	mustWriteFile(t, root, "clean.go", "package main\nimport \"fmt\"\n")
	if err := os.MkdirAll(filepath.Join(root, "sub"), 0755); err != nil {
		t.Fatalf("mkdir sub: %v", err)
	}
	mustWriteFile(t, root, "sub/bad.go", "package sub\nimport \"flag\"\n")

	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatalf("Chdir: %v", err)
	}
	defer os.Chdir(oldDir)

	output, err := captureStdout(t, func() error {
		return handle([]string{"vet", "./..."})
	})
	if err != nil {
		t.Fatalf("handle(vet ./...): %v", err)
	}
	if !strings.Contains(output, "[std-flag]") {
		t.Errorf("expected std-flag violation, got: %s", output)
	}
	if strings.Contains(output, "clean.go") {
		t.Errorf("clean.go should not be reported, got: %s", output)
	}
}

func mustWriteFile(t *testing.T, dir, name, content string) {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("writeFile %s: %v", path, err)
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
