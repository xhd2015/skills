package playwrightdebug

import (
	"os"
	"path/filepath"
)

// DefaultCacheDir returns the default playwright-debug npm package cache location.
func DefaultCacheDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".playwright-debug", "node_package")
}