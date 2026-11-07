package skillcmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"
)

const demoRootContent = "---\nname: demo\ndescription: >-\n  Demo skill with\n  nested topics.\n---\nroot demo\n"

func treeTestRegistry() *Registry {
	tree := fstest.MapFS{
		"a/b/TOPIC.md": &fstest.MapFile{Data: []byte("topic a/b\n")},
		"c/TOPIC.md":   &fstest.MapFile{Data: []byte("topic c\n")},
	}
	return &Registry{
		Skills: []RegisteredSkill{
			{
				Name:    "demo",
				Content: demoRootContent,
				TreeFS:  tree,
			},
			{
				Name:    "flat",
				Content: "no frontmatter here\n",
			},
		},
	}
}

func TestRegistryShowRootAndTopicPath(t *testing.T) {
	reg := treeTestRegistry()
	cases := []struct {
		args []string
		want string
	}{
		{[]string{"--show", "demo"}, demoRootContent},
		{[]string{"--show", "demo/a/b"}, "topic a/b\n"},
		{[]string{"demo/a/b", "--show"}, "topic a/b\n"},
		{[]string{"--show", "demo", "c"}, "topic c\n"},
	}
	for _, tc := range cases {
		out, err := captureStdout(func() error { return reg.HandleSkill(tc.args) })
		if err != nil {
			t.Fatalf("HandleSkill(%v): %v", tc.args, err)
		}
		if out != tc.want {
			t.Fatalf("HandleSkill(%v) = %q, want %q", tc.args, out, tc.want)
		}
	}
}

func TestRegistryShowUnknownTopic(t *testing.T) {
	reg := treeTestRegistry()
	_, err := captureStdout(func() error { return reg.HandleSkill([]string{"--show", "demo/nope"}) })
	if err == nil || !strings.Contains(err.Error(), "nope") {
		t.Fatalf("error = %v, want unknown topic error", err)
	}
}

func TestRegistryShowTopicOnFlatSkill(t *testing.T) {
	reg := treeTestRegistry()
	_, err := captureStdout(func() error { return reg.HandleSkill([]string{"--show", "flat/x"}) })
	if err == nil || !strings.Contains(err.Error(), "unknown topic path: x") {
		t.Fatalf("error = %v, want unknown topic path error", err)
	}
}

func TestRegistryListDescriptionsAndTopics(t *testing.T) {
	reg := treeTestRegistry()
	out, err := captureStdout(func() error { return reg.HandleSkill([]string{"--list"}) })
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	want := "demo - Demo skill with nested topics.\ndemo/a/b\ndemo/c\nflat\n"
	if out != want {
		t.Fatalf("list = %q, want %q", out, want)
	}
}

func TestRegistryDefaultSkill(t *testing.T) {
	reg := treeTestRegistry()
	reg.Default = "demo"
	out, err := captureStdout(func() error { return reg.HandleSkill([]string{"--show"}) })
	if err != nil {
		t.Fatalf("show default: %v", err)
	}
	if out != demoRootContent {
		t.Fatalf("default show = %q", out)
	}
}

func TestRegistryMissingNameWithoutDefault(t *testing.T) {
	reg := treeTestRegistry()
	for _, args := range [][]string{{"--show"}, {"--version"}, {"--install"}} {
		_, err := captureStdout(func() error { return reg.HandleSkill(args) })
		want := "expected skill name for " + args[0]
		if err == nil || !strings.Contains(err.Error(), want) {
			t.Fatalf("HandleSkill(%v) error = %v, want %q", args, err, want)
		}
	}
}

func TestRegistryUnknownSkillNames(t *testing.T) {
	reg := treeTestRegistry()
	_, err := captureStdout(func() error { return reg.HandleSkill([]string{"--show", "nope"}) })
	if err == nil || !strings.Contains(err.Error(), "unknown skill: nope") {
		t.Fatalf("error = %v", err)
	}
	if err == nil || !strings.Contains(err.Error(), "demo, flat") {
		t.Fatalf("error = %v, want available names", err)
	}
}

func TestRegistryInstallTreeSkill(t *testing.T) {
	reg := treeTestRegistry()
	dest := filepath.Join(t.TempDir(), "skills")
	if _, err := captureStdout(func() error {
		return reg.HandleSkill([]string{"--install", "demo", "--dir", dest})
	}); err != nil {
		t.Fatalf("install: %v", err)
	}
	for _, rel := range []string{"SKILL.md", "a/b/TOPIC.md", "c/TOPIC.md"} {
		p := filepath.Join(dest, "demo", filepath.FromSlash(rel))
		if _, err := os.Stat(p); err != nil {
			t.Fatalf("missing installed file %s: %v", rel, err)
		}
	}
	// Topic suffix must not create a nested skill directory of its own.
	if _, err := os.Stat(filepath.Join(dest, "demo", "a")); err != nil {
		t.Fatalf("topic dir missing: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dest, "demo", "a", "SKILL.md")); !os.IsNotExist(err) {
		t.Fatalf("topic installed a SKILL.md, want only TOPIC.md")
	}
}

func TestRegistryInstallTopicSuffixDropsTopic(t *testing.T) {
	reg := treeTestRegistry()
	dest := filepath.Join(t.TempDir(), "skills")
	if _, err := captureStdout(func() error {
		return reg.HandleSkill([]string{"--install", "demo/a/b", "--dir", dest})
	}); err != nil {
		t.Fatalf("install with topic suffix: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dest, "demo", "SKILL.md")); err != nil {
		t.Fatalf("topic-suffixed install did not install skill root: %v", err)
	}
}

func TestRegistryUpdateCollectsTreeFiles(t *testing.T) {
	reg := treeTestRegistry()
	sk, ok := reg.find("demo")
	if !ok {
		t.Fatal("demo not found")
	}
	files, err := reg.single(sk).extraFilesForInstall()
	if err != nil {
		t.Fatalf("extraFilesForInstall: %v", err)
	}
	got := make([]string, 0, len(files))
	for _, f := range files {
		got = append(got, f.Path)
	}
	want := []string{"a/b/TOPIC.md", "c/TOPIC.md"}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("tree files = %v, want %v", got, want)
	}
}

func TestRegistryInstallUsageDefaultAndOverride(t *testing.T) {
	reg := treeTestRegistry()
	if got := reg.installUsage(reg.Skills[0]); got != "skill --install demo" {
		t.Fatalf("default usage = %q", got)
	}
	reg.Usage = "wrk skill --install"
	if got := reg.installUsage(reg.Skills[0]); got != "wrk skill --install demo" {
		t.Fatalf("prefixed usage = %q", got)
	}
	reg.Skills[0].Usage = "wrk skill --install wrk-pipeline"
	if got := reg.installUsage(reg.Skills[0]); got != "wrk skill --install wrk-pipeline" {
		t.Fatalf("explicit usage = %q", got)
	}
}

func TestExtractDescription(t *testing.T) {
	cases := []struct {
		name    string
		content string
		want    string
	}{
		{"folded", "---\nname: d\ndescription: >-\n  One line\n  two.\n---\nbody\n", "One line two."},
		{"literal", "---\ndescription: |\n  Alpha\n  beta\n---\n", "Alpha beta"},
		{"inline", "---\ndescription: Single sentence.\n---\n", "Single sentence."},
		{"quoted", "---\ndescription: \"Quoted: yes\"\n---\n", "Quoted: yes"},
		{"missing", "---\nname: d\n---\n", ""},
		{"no frontmatter", "body only\n", ""},
	}
	for _, tc := range cases {
		if got := ExtractDescription(tc.content); got != tc.want {
			t.Fatalf("%s: ExtractDescription = %q, want %q", tc.name, got, tc.want)
		}
	}
}
