package main

import (
	"io"
	"os"
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
		t.Errorf("missing kool-create in output: %s", output)
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
