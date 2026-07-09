package skillcmd

import (
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
	found := false
	for _, a := range p.Rest {
		if a == "--help" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("Rest missing --help: %v", p.Rest)
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
