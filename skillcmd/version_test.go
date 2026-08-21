package skillcmd

import (
	"errors"
	"strings"
	"testing"
)

func TestParseSkillHeaderMetadata(t *testing.T) {
	content := `---
name: demo
description: >-
  Demo skill.
license: Apache-2.0
compatibility: Requires git
allowed-tools: Bash(git:*) Read
metadata:
  author: example-org
  version: "1.0"
---

# Demo
`
	header, err := ParseSkillHeader(content)
	if err != nil {
		t.Fatalf("ParseSkillHeader: %v", err)
	}
	if header.Name != "demo" || header.Description != "Demo skill." {
		t.Fatalf("unexpected identity: %+v", header)
	}
	if header.License != "Apache-2.0" || header.Compatibility != "Requires git" {
		t.Fatalf("unexpected optional fields: %+v", header)
	}
	if header.AllowedTools != "Bash(git:*) Read" {
		t.Fatalf("AllowedTools = %q", header.AllowedTools)
	}
	if header.Metadata["author"] != "example-org" || header.Metadata["version"] != "1.0" {
		t.Fatalf("unexpected metadata: %#v", header.Metadata)
	}
	version, err := SkillVersion(content)
	if err != nil {
		t.Fatalf("SkillVersion: %v", err)
	}
	if version != "1.0" {
		t.Fatalf("version = %q, want 1.0", version)
	}
}

func TestSkillVersionAcceptsArbitraryString(t *testing.T) {
	version, err := SkillVersion(versionedSkill("demo", "release-candidate"))
	if err != nil {
		t.Fatalf("SkillVersion: %v", err)
	}
	if version != "release-candidate" {
		t.Fatalf("version = %q", version)
	}
}

func TestSkillVersionMissingIsOptionalUntilQueried(t *testing.T) {
	content := "---\nname: demo\ndescription: Demo skill.\n---\nbody\n"
	header, err := ParseSkillHeader(content)
	if err != nil {
		t.Fatalf("ParseSkillHeader: %v", err)
	}
	if len(header.Metadata) != 0 {
		t.Fatalf("metadata = %#v, want empty", header.Metadata)
	}
	_, err = SkillVersion(content)
	if !errors.Is(err, ErrSkillVersionMissing) {
		t.Fatalf("SkillVersion error = %v, want ErrSkillVersionMissing", err)
	}
}

func TestParseSkillHeaderRejectsNonStringMetadata(t *testing.T) {
	for _, tt := range []struct {
		name    string
		content string
		wantErr string
	}{
		{
			name:    "value",
			content: "---\nname: demo\ndescription: Demo skill.\nmetadata:\n  version: 1.0\n---\nbody\n",
			wantErr: "metadata.version must be a string",
		},
		{
			name:    "key",
			content: "---\nname: demo\ndescription: Demo skill.\nmetadata:\n  1: value\n---\nbody\n",
			wantErr: "metadata keys must be strings",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseSkillHeader(tt.content)
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("error = %v, want %q", err, tt.wantErr)
			}
		})
	}
}

func TestParseSkillArgsVersion(t *testing.T) {
	for _, args := range [][]string{{"--version", "demo"}, {"demo", "--version"}} {
		parsed, err := ParseSkillArgs(args)
		if err != nil {
			t.Fatalf("ParseSkillArgs(%v): %v", args, err)
		}
		if parsed.Action != ActionVersion || len(parsed.Rest) != 1 || parsed.Rest[0] != "demo" {
			t.Fatalf("ParseSkillArgs(%v) = %+v", args, parsed)
		}
	}

	parsed, err := ParseSkillArgs([]string{"--version", "--help"})
	if err != nil {
		t.Fatalf("version help: %v", err)
	}
	if parsed.Action != ActionHelp {
		t.Fatalf("version help action = %q, want help", parsed.Action)
	}
}

func TestParseSkillArgsVersionIsExclusive(t *testing.T) {
	for _, args := range [][]string{
		{"--version", "--show"},
		{"--version", "--install"},
		{"--version", "--list"},
		{"--version", "--header"},
	} {
		if _, err := ParseSkillArgs(args); err == nil {
			t.Fatalf("ParseSkillArgs(%v): expected error", args)
		}
	}
}

func TestSingleSkillVersion(t *testing.T) {
	sk := &SingleSkill{
		Name:        "demo",
		RootContent: versionedSkill("demo", "0.1.0"),
	}
	out, err := captureStdout(func() error {
		return sk.Handle([]string{"--version"})
	})
	if err != nil {
		t.Fatalf("Handle --version: %v", err)
	}
	if out != "0.1.0\n" {
		t.Fatalf("stdout = %q, want %q", out, "0.1.0\\n")
	}
	if err := sk.Handle([]string{"--version", "topic"}); err == nil {
		t.Fatal("version with topic: expected error")
	}
}

func TestRegistryVersionBothOrders(t *testing.T) {
	registry := &Registry{Skills: []RegisteredSkill{{
		Name:    "demo",
		Content: versionedSkill("demo", "2.3.4"),
	}}}
	for _, args := range [][]string{{"--version", "demo"}, {"demo", "--version"}} {
		out, err := captureStdout(func() error {
			return registry.HandleSkill(args)
		})
		if err != nil {
			t.Fatalf("HandleSkill(%v): %v", args, err)
		}
		if out != "2.3.4\n" {
			t.Fatalf("HandleSkill(%v) stdout = %q", args, out)
		}
	}
}

func TestRegistryVersionMissing(t *testing.T) {
	registry := &Registry{Skills: []RegisteredSkill{{
		Name:    "demo",
		Content: "---\nname: demo\ndescription: Demo skill.\n---\nbody\n",
	}}}
	_, err := captureStdout(func() error {
		return registry.HandleSkill([]string{"demo", "--version"})
	})
	if err == nil || !strings.Contains(err.Error(), "skill demo has no metadata.version") {
		t.Fatalf("error = %v, want missing version error", err)
	}
}

func versionedSkill(name, version string) string {
	return "---\nname: " + name + "\ndescription: Demo skill.\nmetadata:\n  version: \"" + version + "\"\n---\nbody\n"
}
