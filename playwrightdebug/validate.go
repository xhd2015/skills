package playwrightdebug

import (
	"fmt"
	"os"
	"strings"
)

func hasJSSuffix(path string) bool {
	return strings.HasSuffix(strings.ToLower(path), ".js")
}

// ValidateScriptFile checks that path refers to an existing .js file.
func ValidateScriptFile(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("script file not found: %s", path)
		}
		return fmt.Errorf("cannot access script file %s: %w", path, err)
	}
	if info.IsDir() {
		return fmt.Errorf("run requires an existing .js file: %s is a directory", path)
	}
	if !hasJSSuffix(path) {
		return fmt.Errorf("run requires an existing .js file, not an inline script")
	}
	return nil
}