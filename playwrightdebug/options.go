package playwrightdebug

import "io"

// RunOptions configures a file-mode Playwright script run.
type RunOptions struct {
	ScriptPath        string
	ScriptArgs        []string
	Stdout            io.Writer
	Stderr            io.Writer
	CacheDir          string
	PlaywrightVersion string
	SkipEnsure        bool // test hook: skip npm install

	// Launch controls Chromium startup (default headless vs extension-capable).
	// Zero value keeps legacy headless chromium.launch behavior.
	Launch LaunchOptions
}
