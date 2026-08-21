package skillcmd

import (
	"slices"
	"strings"
	"testing"
)

func TestParseSkillArgsHelpOnly(t *testing.T) {
	for _, args := range [][]string{{"--help"}, {"-h"}} {
		p, err := ParseSkillArgs(args)
		if err != nil {
			t.Fatalf("ParseSkillArgs(%v): %v", args, err)
		}
		if p.Action != ActionHelp {
			t.Fatalf("Action = %q, want help", p.Action)
		}
	}
}

func TestParseSkillArgsShowWithHelpIsSkillHelp(t *testing.T) {
	p, err := ParseSkillArgs([]string{"--show", "--help"})
	if err != nil {
		t.Fatalf("ParseSkillArgs: %v", err)
	}
	if p.Action != ActionHelp {
		t.Fatalf("Action = %q, want help", p.Action)
	}
}

func TestParseSkillArgsInstallWithHelpKeepsInstall(t *testing.T) {
	p, err := ParseSkillArgs([]string{"--install", "--help"})
	if err != nil {
		t.Fatalf("ParseSkillArgs: %v", err)
	}
	if p.Action != ActionInstall {
		t.Fatalf("Action = %q, want install", p.Action)
	}
	if !slices.Contains(p.Rest, "--help") {
		t.Fatalf("Rest missing --help: %v", p.Rest)
	}
}

func TestParseSkillArgsInstallPreservesDownstreamFlags(t *testing.T) {
	for _, args := range [][]string{
		{"--install", "demo", "--global", "target"},
		{"demo", "--install", "--global", "target"},
	} {
		parsed, err := ParseSkillArgs(args)
		if err != nil {
			t.Fatalf("ParseSkillArgs(%v): %v", args, err)
		}
		if parsed.Action != ActionInstall {
			t.Fatalf("ParseSkillArgs(%v) action = %q", args, parsed.Action)
		}
		want := []string{"demo", "--global", "target"}
		if !slices.Equal(parsed.Rest, want) {
			t.Fatalf("ParseSkillArgs(%v) rest = %v, want %v", args, parsed.Rest, want)
		}
	}
}

func TestParseSkillArgsMissingActionHintsHelp(t *testing.T) {
	_, err := ParseSkillArgs([]string{"foo"})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "--help") {
		t.Fatalf("error should mention --help: %v", err)
	}
}
