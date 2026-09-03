package skillcmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveExplicitSkillDir(t *testing.T) {
	t.Parallel()
	skill := "tsk"

	t.Run("skills_basename_always_nests", func(t *testing.T) {
		t.Parallel()
		got, err := ResolveExplicitSkillDir(skill, "/tmp/ai/skills")
		if err != nil {
			t.Fatal(err)
		}
		want := filepath.Join("/tmp/ai/skills", skill)
		if got != want {
			t.Fatalf("got %q want %q", got, want)
		}
	})

	t.Run("skills_basename_wins_over_skill_md", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		collection := filepath.Join(dir, "skills")
		if err := os.MkdirAll(collection, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(collection, "SKILL.md"), []byte("# stray\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		got, err := ResolveExplicitSkillDir(skill, collection)
		if err != nil {
			t.Fatal(err)
		}
		want := filepath.Join(collection, skill)
		if got != want {
			t.Fatalf("got %q want %q", got, want)
		}
	})

	t.Run("existing_skill_md_is_direct", func(t *testing.T) {
		t.Parallel()
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte("# demo\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		got, err := ResolveExplicitSkillDir(skill, dir)
		if err != nil {
			t.Fatal(err)
		}
		if got != dir {
			t.Fatalf("got %q want %q", got, dir)
		}
	})

	t.Run("basename_matches_skill_direct", func(t *testing.T) {
		t.Parallel()
		dir := filepath.Join(t.TempDir(), skill)
		got, err := ResolveExplicitSkillDir(skill, dir)
		if err != nil {
			t.Fatal(err)
		}
		if got != dir {
			t.Fatalf("got %q want %q", got, dir)
		}
	})

	t.Run("parent_nests_skill_name", func(t *testing.T) {
		t.Parallel()
		parent := filepath.Join(t.TempDir(), "out")
		got, err := ResolveExplicitSkillDir(skill, parent)
		if err != nil {
			t.Fatal(err)
		}
		want := filepath.Join(parent, skill)
		if got != want {
			t.Fatalf("got %q want %q", got, want)
		}
	})

	t.Run("empty_dir_errors", func(t *testing.T) {
		t.Parallel()
		if _, err := ResolveExplicitSkillDir(skill, "  "); err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("empty_skill_errors", func(t *testing.T) {
		t.Parallel()
		if _, err := ResolveExplicitSkillDir("", "/tmp/skills"); err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestResolveTargetDirs_DirFlag(t *testing.T) {
	t.Parallel()

	t.Run("dir_and_positional_conflict", func(t *testing.T) {
		t.Parallel()
		_, err := ResolveTargetDirs("tsk", TargetFlags{Dir: "/tmp/skills"}, []string{"/other"})
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("dir_and_cursor_conflict", func(t *testing.T) {
		t.Parallel()
		_, err := ResolveTargetDirs("tsk", TargetFlags{Dir: "/tmp/skills", Cursor: true}, nil)
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("extra_positionals_error", func(t *testing.T) {
		t.Parallel()
		_, err := ResolveTargetDirs("tsk", TargetFlags{}, []string{"/tmp/a", "/tmp/b"})
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("dir_collection_resolves", func(t *testing.T) {
		t.Parallel()
		dirs, err := ResolveTargetDirs("tsk", TargetFlags{Dir: "/tmp/ai/skills"}, nil)
		if err != nil {
			t.Fatal(err)
		}
		if len(dirs) != 1 || dirs[0] != filepath.Join("/tmp/ai/skills", "tsk") {
			t.Fatalf("dirs=%v", dirs)
		}
	})

	t.Run("positional_collection_resolves", func(t *testing.T) {
		t.Parallel()
		dirs, err := ResolveTargetDirs("tsk", TargetFlags{}, []string{"/tmp/ai/skills"})
		if err != nil {
			t.Fatal(err)
		}
		if len(dirs) != 1 || dirs[0] != filepath.Join("/tmp/ai/skills", "tsk") {
			t.Fatalf("dirs=%v", dirs)
		}
	})
}
