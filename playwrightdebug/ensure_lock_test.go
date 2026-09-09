//go:build unix

package playwrightdebug

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestWithEnsureLockSerializesConcurrentAccess(t *testing.T) {
	cacheDir := t.TempDir()
	var active int32
	var maxActive int32
	var wg sync.WaitGroup

	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := withEnsureLock(cacheDir, func() error {
				current := atomic.AddInt32(&active, 1)
				for {
					prev := atomic.LoadInt32(&maxActive)
					if current <= prev || atomic.CompareAndSwapInt32(&maxActive, prev, current) {
						break
					}
				}
				time.Sleep(20 * time.Millisecond)
				atomic.AddInt32(&active, -1)
				return nil
			})
			if err != nil {
				t.Errorf("withEnsureLock: %v", err)
			}
		}()
	}

	wg.Wait()
	if got := atomic.LoadInt32(&maxActive); got != 1 {
		t.Fatalf("max concurrent holders = %d, want 1", got)
	}
}

func TestPlaywrightPackageMatches(t *testing.T) {
	cacheDir := t.TempDir()
	if playwrightPackageMatches(cacheDir, "1.61.0") {
		t.Fatal("expected missing playwright package")
	}

	playwrightDir := filepath.Join(cacheDir, "node_modules", "playwright")
	if err := os.MkdirAll(playwrightDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(playwrightDir, "package.json"), []byte(`{"version":"1.61.0"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if !playwrightPackageMatches(cacheDir, "1.61.0") {
		t.Fatal("expected matching playwright package")
	}
	if playwrightPackageMatches(cacheDir, "1.62.0") {
		t.Fatal("expected mismatched playwright package")
	}
	if !playwrightPackageMatches(cacheDir, "") {
		t.Fatal("expected any valid installed version to match an unspecified version")
	}
}

func TestValidatePlaywrightVersion(t *testing.T) {
	for _, version := range []string{"", "1.61.0", "1.62.0-beta.1"} {
		if err := ValidatePlaywrightVersion(version); err != nil {
			t.Fatalf("ValidatePlaywrightVersion(%q): %v", version, err)
		}
	}
	for _, version := range []string{"latest", "^1.61.0", "file:../../tmp", "1.61"} {
		if err := ValidatePlaywrightVersion(version); err == nil {
			t.Fatalf("ValidatePlaywrightVersion(%q) unexpectedly succeeded", version)
		}
	}
}

func TestVersionedCacheDir(t *testing.T) {
	base := t.TempDir()
	got, err := VersionedCacheDir(base, "1.61.0")
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(base, "versions", "1.61.0")
	if got != want {
		t.Fatalf("VersionedCacheDir = %q, want %q", got, want)
	}
}

const fakeNodeChromium = `#!/bin/sh
# exists check / generic: need chromium-ready
test -f "$PWD/chromium-ready" || exit 1
# path probe for smokeChromiumLoads: print a stub binary
case "$2" in
  *process.stdout.write*)
    stub="$PWD/chromium-stub"
    if [ -f "$PWD/chromium-loads" ]; then
      printf '#!/bin/sh\nexit 0\n' > "$stub"
    else
      printf '#!/bin/sh\necho "error while loading shared libraries: libatk-1.0.so.0" >&2\nexit 127\n' > "$stub"
    fi
    chmod +x "$stub"
    printf '%s' "$stub"
    ;;
esac
`

const fakeNpmEnsure = `#!/bin/sh
echo "npm $*" >> "$COMMAND_LOG"
case "$1" in
  init)
    printf '{"name":"node_package"}\n' > package.json
    ;;
  install)
    target="$3"
    version="${target#playwright@}"
    if [ "$version" = "$target" ]; then version="9.9.9"; fi
    mkdir -p node_modules/playwright
    printf '{"version":"%s"}\n' "$version" > node_modules/playwright/package.json
    ;;
  exec)
    # npm exec -- playwright install|install-deps chromium
    action="$4"
    case "$action" in
      install-deps)
        echo "APT_CONFIG=${APT_CONFIG:-}" >> "$COMMAND_LOG"
        echo "DEBIAN_FRONTEND=${DEBIAN_FRONTEND:-}" >> "$COMMAND_LOG"
        if [ -n "${APT_CONFIG:-}" ] && [ -f "$APT_CONFIG" ]; then
          echo "APT_CONFIG_BODY=$(tr '\n' ' ' < "$APT_CONFIG")" >> "$COMMAND_LOG"
        fi
        if [ "${FAIL_INSTALL_DEPS:-}" = "1" ]; then
          echo "simulated install-deps failure" >&2
          exit 1
        fi
        : > chromium-loads
        ;;
      install)
        : > chromium-ready
        ;;
    esac
    ;;
esac
`

func TestEnsurePlaywrightVersionInstallsAndRepairs(t *testing.T) {
	binDir := t.TempDir()
	logPath := filepath.Join(t.TempDir(), "commands.log")
	writeFakeExecutable(t, filepath.Join(binDir, "node"), fakeNodeChromium)
	writeFakeExecutable(t, filepath.Join(binDir, "npm"), fakeNpmEnsure)
	t.Setenv("COMMAND_LOG", logPath)
	t.Setenv("PATH", binDir+string(os.PathListSeparator)+os.Getenv("PATH"))

	base := t.TempDir()
	var output bytes.Buffer
	cacheDir, err := EnsurePlaywrightVersion(base, "1.61.0", &output, &output)
	if err != nil {
		t.Fatal(err)
	}
	wantCache := filepath.Join(base, "versions", "1.61.0")
	if cacheDir != wantCache {
		t.Fatalf("cacheDir = %q, want %q", cacheDir, wantCache)
	}
	firstLog := readTestFile(t, logPath)
	if !strings.Contains(firstLog, "npm install --save-exact playwright@1.61.0") {
		t.Fatalf("missing exact install command:\n%s", firstLog)
	}
	if !strings.Contains(firstLog, "npm exec -- playwright install chromium") {
		t.Fatalf("missing Chromium install command:\n%s", firstLog)
	}
	if runtime.GOOS == "linux" {
		if !strings.Contains(firstLog, "npm exec -- playwright install-deps chromium") {
			t.Fatalf("linux ensure must install OS deps:\n%s", firstLog)
		}
	} else if strings.Contains(firstLog, "install-deps") {
		t.Fatalf("non-linux ensure must not install-deps:\n%s", firstLog)
	}

	if _, err := EnsurePlaywrightVersion(base, "1.61.0", &output, &output); err != nil {
		t.Fatal(err)
	}
	if secondLog := readTestFile(t, logPath); secondLog != firstLog {
		t.Fatalf("ready cache should skip commands\nfirst:\n%s\nsecond:\n%s", firstLog, secondLog)
	}

	if err := os.Remove(filepath.Join(cacheDir, "chromium-ready")); err != nil {
		t.Fatal(err)
	}
	if _, err := EnsurePlaywrightVersion(base, "1.61.0", &output, &output); err != nil {
		t.Fatal(err)
	}
	repairLog := readTestFile(t, logPath)
	if strings.Count(repairLog, "npm install --save-exact") != 1 {
		t.Fatalf("browser repair should not reinstall matching package:\n%s", repairLog)
	}
	if strings.Count(repairLog, "npm exec -- playwright install chromium") != 2 {
		t.Fatalf("browser repair should rerun Chromium install:\n%s", repairLog)
	}
	// OS deps are system-wide; binary repair should not re-run install-deps when loads still OK.
	if runtime.GOOS == "linux" && strings.Count(repairLog, "npm exec -- playwright install-deps chromium") != 1 {
		t.Fatalf("browser repair should not reinstall OS deps:\n%s", repairLog)
	}
}

func TestEnsurePlaywrightVersionReinstallsMismatchedPackage(t *testing.T) {
	binDir := t.TempDir()
	logPath := filepath.Join(t.TempDir(), "commands.log")
	writeFakeExecutable(t, filepath.Join(binDir, "node"), fakeNodeChromium)
	writeFakeExecutable(t, filepath.Join(binDir, "npm"), fakeNpmEnsure)
	t.Setenv("COMMAND_LOG", logPath)
	t.Setenv("PATH", binDir+string(os.PathListSeparator)+os.Getenv("PATH"))

	base := t.TempDir()
	cacheDir := filepath.Join(base, "versions", "1.61.0")
	if err := os.MkdirAll(filepath.Join(cacheDir, "node_modules", "playwright"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(cacheDir, "node_modules", "playwright", "package.json"), []byte(`{"version":"1.62.0"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(cacheDir, "chromium-ready"), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(cacheDir, "chromium-loads"), nil, 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := EnsurePlaywrightVersion(base, "1.61.0", nil, nil); err != nil {
		t.Fatal(err)
	}
	if got := readTestFile(t, logPath); !strings.Contains(got, "npm install --save-exact playwright@1.61.0") {
		t.Fatalf("mismatched package was not reinstalled:\n%s", got)
	}
}

func TestEnsurePlaywrightVersionInstallsOSDepsWhenBinaryPresent(t *testing.T) {
	prevNeeds := needsChromiumOSDeps
	needsChromiumOSDeps = func() bool { return true }
	t.Cleanup(func() { needsChromiumOSDeps = prevNeeds })
	prevHost := lookupHostOS
	lookupHostOS = func() HostOS {
		return HostOS{ID: "debian", VersionID: "11", VersionCodename: "bullseye", Arch: "amd64"}
	}
	t.Cleanup(func() { lookupHostOS = prevHost })

	binDir := t.TempDir()
	logPath := filepath.Join(t.TempDir(), "commands.log")
	writeFakeExecutable(t, filepath.Join(binDir, "node"), fakeNodeChromium)
	writeFakeExecutable(t, filepath.Join(binDir, "npm"), fakeNpmEnsure)
	t.Setenv("COMMAND_LOG", logPath)
	t.Setenv("PATH", binDir+string(os.PathListSeparator)+os.Getenv("PATH"))

	base := t.TempDir()
	cacheDir := filepath.Join(base, "versions", "1.61.0")
	if err := os.MkdirAll(filepath.Join(cacheDir, "node_modules", "playwright"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(cacheDir, "node_modules", "playwright", "package.json"), []byte(`{"version":"1.61.0"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(cacheDir, "package.json"), []byte(`{"name":"node_package"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	// Binary on disk, OS deps missing (no chromium-loads) — the debian11 failure mode.
	if err := os.WriteFile(filepath.Join(cacheDir, "chromium-ready"), nil, 0o644); err != nil {
		t.Fatal(err)
	}

	var output bytes.Buffer
	if _, err := EnsurePlaywrightVersion(base, "1.61.0", &output, &output); err != nil {
		t.Fatal(err)
	}
	got := readTestFile(t, logPath)
	if strings.Contains(got, "npm exec -- playwright install chromium") {
		t.Fatalf("should not reinstall Chromium when binary exists:\n%s", got)
	}
	if !strings.Contains(got, "npm exec -- playwright install-deps chromium") {
		t.Fatalf("missing install-deps repair:\n%s", got)
	}
	if !strings.Contains(got, "APT_CONFIG=") || !strings.Contains(got, "Check-Valid-Until") {
		t.Fatalf("install-deps must set APT_CONFIG ignoring expired Release:\n%s", got)
	}
	if !strings.Contains(got, "DEBIAN_FRONTEND=noninteractive") {
		t.Fatalf("install-deps must set DEBIAN_FRONTEND=noninteractive:\n%s", got)
	}
	if !strings.Contains(got, "Dir::Etc::sourcelist") || !strings.Contains(got, "Dir::Etc::sourceparts") {
		t.Fatalf("install-deps on debian must override apt sourcelist:\n%s", got)
	}
	aptConf := filepath.Join(cacheDir, ".playwright-debug-apt.conf")
	if _, err := os.Stat(aptConf); err != nil {
		t.Fatalf("expected apt conf at %s: %v", aptConf, err)
	}
	sources := filepath.Join(cacheDir, ".playwright-debug-apt-sources.list")
	srcBody, err := os.ReadFile(sources)
	if err != nil {
		t.Fatalf("expected main-only sources at %s: %v", sources, err)
	}
	if !strings.Contains(string(srcBody), "bullseye main") || strings.Contains(string(srcBody), "security") {
		t.Fatalf("main-only sources = %q", srcBody)
	}
	if _, err := os.Stat(filepath.Join(cacheDir, "chromium-loads")); err != nil {
		t.Fatalf("install-deps should mark deps ready: %v", err)
	}

	// Second ensure is a no-op.
	before := got
	if _, err := EnsurePlaywrightVersion(base, "1.61.0", &output, &output); err != nil {
		t.Fatal(err)
	}
	if after := readTestFile(t, logPath); after != before {
		t.Fatalf("ready after deps should skip commands\nbefore:\n%s\nafter:\n%s", before, after)
	}
}

func TestEnsurePlaywrightVersionInstallDepsFailureIsExplicit(t *testing.T) {
	prevNeeds := needsChromiumOSDeps
	needsChromiumOSDeps = func() bool { return true }
	t.Cleanup(func() { needsChromiumOSDeps = prevNeeds })

	binDir := t.TempDir()
	logPath := filepath.Join(t.TempDir(), "commands.log")
	writeFakeExecutable(t, filepath.Join(binDir, "node"), fakeNodeChromium)
	writeFakeExecutable(t, filepath.Join(binDir, "npm"), fakeNpmEnsure)
	t.Setenv("COMMAND_LOG", logPath)
	t.Setenv("FAIL_INSTALL_DEPS", "1")
	t.Setenv("PATH", binDir+string(os.PathListSeparator)+os.Getenv("PATH"))

	base := t.TempDir()
	cacheDir := filepath.Join(base, "versions", "1.61.0")
	if err := os.MkdirAll(filepath.Join(cacheDir, "node_modules", "playwright"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(cacheDir, "node_modules", "playwright", "package.json"), []byte(`{"version":"1.61.0"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(cacheDir, "package.json"), []byte(`{"name":"node_package"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(cacheDir, "chromium-ready"), nil, 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := EnsurePlaywrightVersion(base, "1.61.0", nil, nil)
	if err == nil {
		t.Fatal("expected install-deps failure")
	}
	msg := err.Error()
	if !strings.Contains(msg, "install-deps") {
		t.Fatalf("error should mention install-deps:\n%s", msg)
	}
	if !strings.Contains(msg, "security.debian.org") && !strings.Contains(msg, "Check-Valid-Until") {
		t.Fatalf("error should mention apt/security workaround:\n%s", msg)
	}
}

func TestInstallDepsEnvSetsAptConfig(t *testing.T) {
	prevHost := lookupHostOS
	lookupHostOS = func() HostOS {
		return HostOS{ID: "debian", VersionID: "11", VersionCodename: "bullseye", Arch: "amd64"}
	}
	t.Cleanup(func() { lookupHostOS = prevHost })

	cacheDir := t.TempDir()
	env, err := installDepsEnv(cacheDir)
	if err != nil {
		t.Fatal(err)
	}
	var aptConfig, frontend string
	for _, e := range env {
		if strings.HasPrefix(e, "APT_CONFIG=") {
			aptConfig = strings.TrimPrefix(e, "APT_CONFIG=")
		}
		if strings.HasPrefix(e, "DEBIAN_FRONTEND=") {
			frontend = strings.TrimPrefix(e, "DEBIAN_FRONTEND=")
		}
	}
	if aptConfig == "" {
		t.Fatal("missing APT_CONFIG")
	}
	if frontend != "noninteractive" {
		t.Fatalf("DEBIAN_FRONTEND=%q, want noninteractive", frontend)
	}
	body, err := os.ReadFile(aptConfig)
	if err != nil {
		t.Fatal(err)
	}
	text := string(body)
	if !strings.Contains(text, "Acquire::Check-Valid-Until") || !strings.Contains(text, "false") {
		t.Fatalf("apt conf body = %q", body)
	}
	if !strings.Contains(text, "Dir::Etc::sourcelist") || !strings.Contains(text, "Dir::Etc::sourceparts") {
		t.Fatalf("apt conf missing sourcelist override: %q", body)
	}
	src, err := os.ReadFile(aptMainOnlySourcesPath(cacheDir))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(src), "deb http://deb.debian.org/debian bullseye main") {
		t.Fatalf("sources = %q", src)
	}
	if strings.Contains(string(src), "security") {
		t.Fatalf("sources must omit security: %q", src)
	}
}

func writeFakeExecutable(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o755); err != nil {
		t.Fatal(err)
	}
}

func readTestFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}
