package playwrightdebug

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Env keys passed into the Node bootstrap / eval wrapper for launch mode.
const (
	EnvExtensionPaths = "PLAYWRIGHT_DEBUG_EXTENSION_PATHS" // os.PathListSeparator-joined abs paths
	EnvUserDataDir    = "PLAYWRIGHT_DEBUG_USER_DATA_DIR"
	EnvHeaded         = "PLAYWRIGHT_DEBUG_HEADED" // "1" or "0"
	EnvLaunchMode     = "PLAYWRIGHT_DEBUG_LAUNCH_MODE" // "default" | "extension"
)

// LaunchOptions controls how Chromium is started for file/eval runs.
type LaunchOptions struct {
	// ExtensionPaths are unpacked extension directories (must contain manifest.json).
	// When non-empty, bootstrap uses launchPersistentContext + load-extension.
	ExtensionPaths []string
	// UserDataDir is the Chromium profile directory for launchPersistentContext.
	// When set (with or without extensions), cookies/localStorage/login survive
	// process restarts and reboots as long as the path is durable (not under /tmp).
	// The profile is never auto-deleted.
	// Empty + extension mode → Node creates a temp profile and removes it on close.
	// Empty + default mode → ephemeral chromium.launch (no profile).
	UserDataDir string
	// Headed forces a visible browser. When nil:
	//   - extension mode defaults to headed (required for reliable extension load)
	//   - default mode stays headless
	Headed *bool
}

// HasExtension returns true when one or more extension paths are set.
func (o LaunchOptions) HasExtension() bool {
	return len(o.ExtensionPaths) > 0
}

// Normalize validates and absolutizes paths.
func (o LaunchOptions) Normalize() (LaunchOptions, error) {
	out := o
	var paths []string
	for _, p := range o.ExtensionPaths {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		abs, err := filepath.Abs(p)
		if err != nil {
			return out, fmt.Errorf("extension path: %w", err)
		}
		st, err := os.Stat(abs)
		if err != nil {
			return out, fmt.Errorf("extension path %q: %w", abs, err)
		}
		if !st.IsDir() {
			return out, fmt.Errorf("extension path %q is not a directory", abs)
		}
		manifest := filepath.Join(abs, "manifest.json")
		if _, err := os.Stat(manifest); err != nil {
			return out, fmt.Errorf("extension path %q: missing manifest.json", abs)
		}
		paths = append(paths, abs)
	}
	out.ExtensionPaths = paths
	if strings.TrimSpace(o.UserDataDir) != "" {
		abs, err := filepath.Abs(o.UserDataDir)
		if err != nil {
			return out, fmt.Errorf("user-data-dir: %w", err)
		}
		out.UserDataDir = abs
	}
	return out, nil
}

// ApplyEnv merges LaunchOptions into a process env slice (typically PlaywrightEnv output).
func (o LaunchOptions) ApplyEnv(env []string) []string {
	// Strip prior keys so re-runs are clean.
	env = stripEnvKeys(env, EnvExtensionPaths, EnvUserDataDir, EnvHeaded, EnvLaunchMode)

	if o.HasExtension() {
		env = append(env, EnvLaunchMode+"=extension")
		env = append(env, EnvExtensionPaths+"="+strings.Join(o.ExtensionPaths, string(os.PathListSeparator)))
	} else {
		env = append(env, EnvLaunchMode+"=default")
	}
	if o.UserDataDir != "" {
		env = append(env, EnvUserDataDir+"="+o.UserDataDir)
	}
	headed := false
	if o.Headed != nil {
		headed = *o.Headed
	} else if o.HasExtension() {
		headed = true // reliable extension load
	}
	if headed {
		env = append(env, EnvHeaded+"=1")
	} else {
		env = append(env, EnvHeaded+"=0")
	}
	return env
}

func stripEnvKeys(env []string, keys ...string) []string {
	drop := map[string]bool{}
	for _, k := range keys {
		drop[k] = true
	}
	out := env[:0]
	for _, e := range env {
		key := e
		if i := strings.IndexByte(e, '='); i >= 0 {
			key = e[:i]
		}
		if drop[key] {
			continue
		}
		out = append(out, e)
	}
	return out
}

// ExtractLaunchFlags peels playwright-debug launch flags from argv.
// Remaining args are returned for command routing / script args.
//
// Supported flags (may appear anywhere except inside -e/--eval script payload
// after the script string — callers should peel before splitting eval body):
//
//	--extension <dir>       (repeatable)
//	--load-extension <dir>  (alias)
//	--user-data-dir <dir>
//	--headed
//	--headless
func ExtractLaunchFlags(argv []string) (LaunchOptions, []string, error) {
	var opts LaunchOptions
	rest := make([]string, 0, len(argv))
	headedSet := false
	headed := false

	for i := 0; i < len(argv); i++ {
		a := argv[i]
		switch a {
		case "--extension", "--load-extension":
			if i+1 >= len(argv) {
				return opts, nil, fmt.Errorf("%s requires a directory argument", a)
			}
			i++
			opts.ExtensionPaths = append(opts.ExtensionPaths, argv[i])
		case "--user-data-dir":
			if i+1 >= len(argv) {
				return opts, nil, fmt.Errorf("--user-data-dir requires a directory argument")
			}
			i++
			opts.UserDataDir = argv[i]
		case "--headed":
			headedSet = true
			headed = true
		case "--headless":
			headedSet = true
			headed = false
		default:
			// also support --extension=path form
			if strings.HasPrefix(a, "--extension=") {
				opts.ExtensionPaths = append(opts.ExtensionPaths, strings.TrimPrefix(a, "--extension="))
				continue
			}
			if strings.HasPrefix(a, "--load-extension=") {
				opts.ExtensionPaths = append(opts.ExtensionPaths, strings.TrimPrefix(a, "--load-extension="))
				continue
			}
			if strings.HasPrefix(a, "--user-data-dir=") {
				opts.UserDataDir = strings.TrimPrefix(a, "--user-data-dir=")
				continue
			}
			rest = append(rest, a)
		}
	}
	if headedSet {
		opts.Headed = &headed
	}
	norm, err := opts.Normalize()
	if err != nil {
		return opts, rest, err
	}
	return norm, rest, nil
}
