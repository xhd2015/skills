package playwrightdebug

import "io"

// RunOptions configures a file-mode Playwright script run.
type RunOptions struct {
	ScriptPath string
	ScriptArgs []string
	Stdout     io.Writer
	Stderr     io.Writer
	CacheDir   string
	SkipEnsure bool // test hook: skip npm install
}