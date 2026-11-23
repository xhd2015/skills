package skillcmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// LinkResult is what LinkToDir did, or what it found already true.
//
// Three outcomes rather than a bool, because the caller has to say something
// different in each case: it created a link, there was nothing to do, or
// somebody else's link is in the way. A bool would collapse the last two, and
// "skipped" is the wrong thing to print about a link that points somewhere
// unexpected.
type LinkResult string

const (
	// LinkCreated means the entry did not exist and now points at src.
	LinkCreated LinkResult = "created"
	// LinkPresent means the entry already pointed at src.
	LinkPresent LinkResult = "present"
	// LinkConflict means the entry exists and points elsewhere -- or is a real
	// directory, or a file. It is never replaced: a link this function did not
	// create is somebody's decision, and clobbering it silently is how a global
	// skills directory loses work.
	LinkConflict LinkResult = "conflict"
)

// LinkToDir makes destDir/<basename(src)> a symbolic link to src.
//
// The one place every "put this skill in the global directory" flow goes
// through, so a second caller cannot invent a fourth outcome or a different
// answer for a link that already points at the wrong place.
//
// Idempotent: running it twice returns LinkPresent, not a second link. Paths are
// compared after resolving, because the two sides arrive in different shapes --
// a caller may pass the link or the target, an absolute path or one with ".." in
// it -- and a string comparison would report a conflict where the filesystem sees
// the same file.
func LinkToDir(src, destDir string) (LinkResult, error) {
	if src == "" || destDir == "" {
		return "", fmt.Errorf("skillcmd: link needs both a source and a destination")
	}
	absSrc, err := filepath.Abs(src)
	if err != nil {
		return "", fmt.Errorf("resolve link source %s: %w", src, err)
	}
	if info, err := os.Stat(absSrc); err != nil || !info.IsDir() {
		return "", fmt.Errorf("link source is not a directory: %s", absSrc)
	}
	absDest, err := filepath.Abs(destDir)
	if err != nil {
		return "", fmt.Errorf("resolve link destination %s: %w", destDir, err)
	}
	link := filepath.Join(absDest, filepath.Base(absSrc))

	// Readlink, not Stat: a symlink to a missing directory still has to be
	// reported as a conflict rather than silently replaced.
	if current, err := os.Readlink(link); err == nil {
		if samePath(current, link, absSrc) {
			return LinkPresent, nil
		}
		// Pointing elsewhere. Left alone on purpose; see LinkConflict.
		return LinkConflict, nil
	} else if !errors.Is(err, os.ErrNotExist) {
		// Not a symlink. It can still be a directory or a file that exists.
		return LinkConflict, nil
	}
	if err := os.MkdirAll(absDest, 0o755); err != nil {
		return "", fmt.Errorf("create %s: %w", absDest, err)
	}
	if err := os.Symlink(absSrc, link); err != nil {
		// Lost a race with another process that created the entry between the
		// Readlink above and here. Its link wins; ours is not needed.
		if current, readErr := os.Readlink(link); readErr == nil && samePath(current, link, absSrc) {
			return LinkPresent, nil
		}
		return "", fmt.Errorf("link %s: %w", link, err)
	}
	return LinkCreated, nil
}

// samePath reports whether two paths name the same directory.
//
// A relative readlink is resolved against the link's own directory, which is what
// the kernel does; resolving against the process cwd is the classic way to think
// a perfectly good link is broken.
func samePath(a, link, want string) bool {
	if !filepath.IsAbs(a) {
		a = filepath.Join(filepath.Dir(link), a)
	}
	resolvedA, err := filepath.EvalSymlinks(a)
	if err != nil {
		// A dangling target: compare the cleaned paths, which is the most that
		// can be said about it.
		return filepath.Clean(a) == filepath.Clean(want)
	}
	resolvedWant, err := filepath.EvalSymlinks(want)
	if err != nil {
		resolvedWant = filepath.Clean(want)
	}
	return resolvedA == resolvedWant
}
