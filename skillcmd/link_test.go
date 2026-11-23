package skillcmd

import (
	"os"
	"path/filepath"
	"testing"
)

// The link is the step that puts a skill where every agent can find it, so its
// two failure modes matter more than its success: replacing somebody else's link,
// and reporting a perfectly good link as broken because it was written relative.

func TestLinkToDirCreatesThenReportsPresent(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "skills", "do-verification")
	dest := filepath.Join(root, "global")
	if err := os.MkdirAll(src, 0o755); err != nil {
		t.Fatal(err)
	}

	if got, err := LinkToDir(src, dest); err != nil || got != LinkCreated {
		t.Fatalf("first link = %q, %v", got, err)
	}
	// Creating the destination is part of the job: a fresh machine has no
	// ~/.agents at all.
	link := filepath.Join(dest, "do-verification")
	if target, err := os.Readlink(link); err != nil || target != src {
		t.Fatalf("link = %q, %v", target, err)
	}

	// Idempotent. A second run must not report work it did not do.
	if got, err := LinkToDir(src, dest); err != nil || got != LinkPresent {
		t.Fatalf("second link = %q, %v", got, err)
	}
	// And it must not have replaced the link with a fresh one.
	if target, err := os.Readlink(link); err != nil || target != src {
		t.Fatalf("relinked = %q, %v", target, err)
	}
}

func TestLinkToDirLeavesAConflictingLinkAlone(t *testing.T) {
	root := t.TempDir()
	mine := filepath.Join(root, "skills", "do-verification")
	theirs := filepath.Join(root, "elsewhere", "do-verification")
	dest := filepath.Join(root, "global")
	for _, dir := range []string{mine, theirs, dest} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	link := filepath.Join(dest, "do-verification")
	if err := os.Symlink(theirs, link); err != nil {
		t.Fatal(err)
	}

	if got, err := LinkToDir(mine, dest); err != nil || got != LinkConflict {
		t.Fatalf("conflict = %q, %v", got, err)
	}
	// The whole point: somebody else's link survives.
	if target, err := os.Readlink(link); err != nil || target != theirs {
		t.Fatalf("the conflicting link was rewritten to %q, %v", target, err)
	}
}

// A real directory in the way is a conflict too. Replacing it would delete files.
func TestLinkToDirDoesNotReplaceADirectory(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "skills", "do-verification")
	dest := filepath.Join(root, "global")
	occupant := filepath.Join(dest, "do-verification")
	for _, dir := range []string{src, occupant} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	kept := filepath.Join(occupant, "SKILL.md")
	if err := os.WriteFile(kept, []byte("hand written"), 0o644); err != nil {
		t.Fatal(err)
	}

	if got, err := LinkToDir(src, dest); err != nil || got != LinkConflict {
		t.Fatalf("directory in the way = %q, %v", got, err)
	}
	if body, err := os.ReadFile(kept); err != nil || string(body) != "hand written" {
		t.Fatalf("the occupant's file was touched: %q, %v", string(body), err)
	}
}

// A relative symlink is what "ln -s ../skills/x" produces, and it is not broken.
func TestLinkToDirReadsARelativeLink(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "skills", "do-verification")
	dest := filepath.Join(root, "agents", "skills")
	for _, dir := range []string{src, dest} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	rel, err := filepath.Rel(dest, src)
	if err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(dest, "do-verification")
	if err := os.Symlink(rel, link); err != nil {
		t.Fatal(err)
	}

	if got, err := LinkToDir(src, dest); err != nil || got != LinkPresent {
		t.Fatalf("a relative link to the right place = %q, %v", got, err)
	}
	if got, err := LinkToDir(src+"/./", dest); err != nil || got != LinkPresent {
		t.Fatalf("an unnormalised source = %q, %v", got, err)
	}
}

func TestLinkToDirRefusesAMissingSource(t *testing.T) {
	root := t.TempDir()
	if _, err := LinkToDir(filepath.Join(root, "absent"), filepath.Join(root, "dest")); err == nil {
		t.Fatal("linked a directory that does not exist")
	}
}
