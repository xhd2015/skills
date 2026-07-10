package install

import (
	"os"
	"path/filepath"
	"testing"
)

func TestHandleInstallGlobalCursor(t *testing.T) {
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)

	err := HandleInstall(InstallOptions{
		SkillDirName: "example-skill",
		SkillContent: "# test skill\n",
	}, []string{"--cursor", "--global"})
	if err != nil {
		t.Fatalf("HandleInstall(--cursor --global): %v", err)
	}

	skillFile := filepath.Join(homeDir, ".cursor", "skills", "example-skill", "SKILL.md")
	content, err := os.ReadFile(skillFile)
	if err != nil {
		t.Fatalf("read skill file: %v", err)
	}
	if string(content) != "# test skill\n" {
		t.Fatalf("unexpected skill content: %q", string(content))
	}
}

func TestHandleInstallGlobalCodex(t *testing.T) {
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)

	err := HandleInstall(InstallOptions{
		SkillDirName: "example-skill",
		SkillContent: "# test skill\n",
	}, []string{"--codex", "--global"})
	if err != nil {
		t.Fatalf("HandleInstall(--codex --global): %v", err)
	}

	skillFile := filepath.Join(homeDir, ".codex", "skills", "example-skill", "SKILL.md")
	content, err := os.ReadFile(skillFile)
	if err != nil {
		t.Fatalf("read skill file: %v", err)
	}
	if string(content) != "# test skill\n" {
		t.Fatalf("unexpected skill content: %q", string(content))
	}
}

func TestHandleInstallGlobalDefaultsToHomeGeneralAgents(t *testing.T) {
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)

	err := HandleInstall(InstallOptions{
		SkillDirName: "example-skill",
		SkillContent: "# test skill\n",
	}, []string{"--global"})
	if err != nil {
		t.Fatalf("HandleInstall(--global): %v", err)
	}

	skillFile := filepath.Join(homeDir, ".agents", "skills", "example-skill", "SKILL.md")
	content, err := os.ReadFile(skillFile)
	if err != nil {
		t.Fatalf("read skill file: %v", err)
	}
	if string(content) != "# test skill\n" {
		t.Fatalf("unexpected skill content: %q", string(content))
	}
}

func TestHandleInstallGlobalCustomDir(t *testing.T) {
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)

	err := HandleInstall(InstallOptions{
		SkillDirName: "example-skill",
		SkillContent: "# test skill\n",
	}, []string{"--global", "my-custom-skill"})
	if err != nil {
		t.Fatalf("HandleInstall(--global custom-dir): %v", err)
	}

	skillFile := filepath.Join(homeDir, "my-custom-skill", "SKILL.md")
	content, err := os.ReadFile(skillFile)
	if err != nil {
		t.Fatalf("read skill file: %v", err)
	}
	if string(content) != "# test skill\n" {
		t.Fatalf("unexpected skill content: %q", string(content))
	}
}

func TestHandleInstallGlobalAbsoluteDirUnchanged(t *testing.T) {
	tmpDir := t.TempDir()
	homeDir := t.TempDir()
	t.Setenv("HOME", homeDir)

	targetDir := filepath.Join(tmpDir, "absolute-target")

	err := HandleInstall(InstallOptions{
		SkillDirName: "example-skill",
		SkillContent: "# test skill\n",
	}, []string{"--global", targetDir})
	if err != nil {
		t.Fatalf("HandleInstall(--global absolute-dir): %v", err)
	}

	skillFile := filepath.Join(targetDir, "SKILL.md")
	content, err := os.ReadFile(skillFile)
	if err != nil {
		t.Fatalf("read skill file: %v", err)
	}
	if string(content) != "# test skill\n" {
		t.Fatalf("unexpected skill content: %q", string(content))
	}
}

