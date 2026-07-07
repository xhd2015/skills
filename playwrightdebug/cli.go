package playwrightdebug

import (
	"fmt"
	"os"
)

// ParseCLIFileRun extracts script path and args from CLI argv for file-mode runs.
func ParseCLIFileRun(argv []string) (scriptPath string, scriptArgs []string, err error) {
	if len(argv) == 0 {
		return "", nil, fmt.Errorf("run requires a .js script file argument")
	}

	switch argv[0] {
	case "run":
		if len(argv) < 2 {
			return "", nil, fmt.Errorf("run requires a .js script file argument")
		}
		return argv[1], argv[2:], nil
	default:
		if isScriptFile(argv[0]) {
			return argv[0], argv[1:], nil
		}
		return "", nil, fmt.Errorf("run requires an existing .js file, not an inline script")
	}
}

func isScriptFile(path string) bool {
	if !hasJSSuffix(path) {
		return false
	}
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}