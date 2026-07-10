package install

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHandleInstallDeletesOrphanWhenPlanMatches(t *testing.T) {
	tmpDir := t.TempDir()
	targetDir := filepath.Join(tmpDir, "example-skill")
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		t.Fatalf("mkdir target: %v", err)
	}
	skillFile := filepath.Join(targetDir, "SKILL.md")
	if err := os.WriteFile(skillFile, []byte("# same skill\n"), 0644); err != nil {
		t.Fatalf("write skill file: %v", err)
	}
	staleFile := filepath.Join(targetDir, "stale.txt")
	if err := os.WriteFile(staleFile, []byte("stale"), 0644); err != nil {
		t.Fatalf("write stale file: %v", err)
	}

	stdout, err := captureStdout(t, func() error {
		return HandleInstall(InstallOptions{
			SkillDirName: "example-skill",
			SkillContent: "# same skill\n",
		}, []string{targetDir})
	})
	if err != nil {
		t.Fatalf("HandleInstall: %v\nstdout:\n%s", err, stdout)
	}
	if strings.Contains(stdout, "Skill is up to date") {
		t.Fatalf("orphan must not be treated as up to date:\n%s", stdout)
	}
	if !strings.Contains(stdout, "Update skill at") {
		t.Fatalf("stdout missing Update header:\n%s", stdout)
	}
	if !strings.Contains(stdout, "delete: "+staleFile) {
		t.Fatalf("stdout missing delete for orphan:\n%s", stdout)
	}
	if _, err := os.Stat(staleFile); !os.IsNotExist(err) {
		t.Fatalf("orphan should be deleted; stat err=%v", err)
	}
	content, err := os.ReadFile(skillFile)
	if err != nil {
		t.Fatalf("read skill file: %v", err)
	}
	if string(content) != "# same skill\n" {
		t.Fatalf("unexpected skill content: %q", string(content))
	}
}
