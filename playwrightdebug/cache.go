package playwrightdebug

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"golang.org/x/mod/semver"
)

var exactPlaywrightVersionPattern = regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+(?:-[0-9A-Za-z]+(?:[.-][0-9A-Za-z]+)*)?$`)

// DefaultCacheDir returns the default playwright-debug npm package cache location.
func DefaultCacheDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".playwright-debug", "node_package")
}

// ValidatePlaywrightVersion accepts only an exact semantic version. Keeping
// arbitrary npm package specs out of this value also makes it safe to use as a
// cache directory name.
func ValidatePlaywrightVersion(version string) error {
	if version == "" {
		return nil
	}
	if !exactPlaywrightVersionPattern.MatchString(version) || !semver.IsValid("v"+version) {
		return fmt.Errorf("invalid Playwright version %q: expected exact version such as 1.61.0", version)
	}
	return nil
}

// VersionedCacheDir isolates explicitly requested Playwright versions while
// preserving the legacy cache location when no version is specified.
func VersionedCacheDir(baseDir, version string) (string, error) {
	if err := ValidatePlaywrightVersion(version); err != nil {
		return "", err
	}
	if baseDir == "" {
		baseDir = DefaultCacheDir()
	}
	if version == "" {
		return baseDir, nil
	}
	return filepath.Join(baseDir, "versions", version), nil
}
