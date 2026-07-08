//go:build windows

package playwrightdebug

import (
	"fmt"
	"os"
)

func withEnsureLock(cacheDir string, fn func() error) error {
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		return fmt.Errorf("create cache dir: %w", err)
	}
	return fn()
}