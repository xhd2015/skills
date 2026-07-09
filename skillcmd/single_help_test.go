package skillcmd

import (
	"io"
	"os"
	"strings"
	"testing"
	"testing/fstest"
)

func TestSingleSkillHandleHelp(t *testing.T) {
	sk := &SingleSkill{
		Name:  "demo",
		Usage: "demo skill --install",
		Help:  "Usage: demo skill --show\n",
	}
	out, err := captureStdout(func() error {
		return sk.Handle([]string{"--help"})
	})
	if err != nil {
		t.Fatalf("Handle: %v", err)
	}
	if !strings.Contains(out, "Usage: demo skill --show") {
		t.Fatalf("stdout missing custom help:\n%s", out)
	}
}

func TestSingleSkillHandleDefaultHelp(t *testing.T) {
	sk := &SingleSkill{Name: "demo", Usage: "demo skill --install"}
	out, err := captureStdout(func() error {
		return sk.Handle([]string{"-h"})
	})
	if err != nil {
		t.Fatalf("Handle: %v", err)
	}
	for _, want := range []string{"--show", "--install", "--list", "demo"} {
		if !strings.Contains(out, want) {
			t.Fatalf("default help missing %q:\n%s", want, out)
		}
	}
}

func TestSingleSkillListAndHelpIncludeTopics(t *testing.T) {
	tree := fstest.MapFS{
		"flags-parsing/SKILL.md":       &fstest.MapFile{Data: []byte("fp\n")},
		"flags-parsing/types/SKILL.md": &fstest.MapFile{Data: []byte("types\n")},
		"skill-cli/SKILL.md":           &fstest.MapFile{Data: []byte("cli\n")},
	}
	sk := &SingleSkill{
		Name:        "demo",
		RootContent: "---\nname: demo\n---\nroot\n",
		TreeFS:      tree,
		Usage:       "demo skill --install",
		Help:        "Usage: demo skill --show\n",
	}

	listOut, err := captureStdout(func() error {
		return sk.Handle([]string{"--list"})
	})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	for _, want := range []string{"demo", "flags-parsing", "flags-parsing/types", "skill-cli"} {
		if !strings.Contains(listOut, want) {
			t.Fatalf("list missing %q:\n%s", want, listOut)
		}
	}

	helpOut, err := captureStdout(func() error {
		return sk.Handle([]string{"--help"})
	})
	if err != nil {
		t.Fatalf("help: %v", err)
	}
	for _, want := range []string{"Usage: demo skill --show", "Available topics:", "flags-parsing", "types"} {
		if !strings.Contains(helpOut, want) {
			t.Fatalf("help missing %q:\n%s", want, helpOut)
		}
	}
}

func captureStdout(fn func() error) (string, error) {
	r, w, err := os.Pipe()
	if err != nil {
		return "", err
	}
	old := os.Stdout
	os.Stdout = w
	runErr := fn()
	_ = w.Close()
	os.Stdout = old
	data, readErr := io.ReadAll(r)
	_ = r.Close()
	if readErr != nil {
		return "", readErr
	}
	return string(data), runErr
}
