//go:build unix

package playwrightdebug

import (
	"bytes"
	"os"
	"path/filepath"
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

func TestEnsurePlaywrightVersionInstallsAndRepairs(t *testing.T) {
	binDir := t.TempDir()
	logPath := filepath.Join(t.TempDir(), "commands.log")
	writeFakeExecutable(t, filepath.Join(binDir, "node"), `#!/bin/sh
test -f "$PWD/chromium-ready"
`)
	writeFakeExecutable(t, filepath.Join(binDir, "npm"), `#!/bin/sh
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
    : > chromium-ready
    ;;
esac
`)
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
}

func TestEnsurePlaywrightVersionReinstallsMismatchedPackage(t *testing.T) {
	binDir := t.TempDir()
	logPath := filepath.Join(t.TempDir(), "commands.log")
	writeFakeExecutable(t, filepath.Join(binDir, "node"), `#!/bin/sh
test -f "$PWD/chromium-ready"
`)
	writeFakeExecutable(t, filepath.Join(binDir, "npm"), `#!/bin/sh
echo "npm $*" >> "$COMMAND_LOG"
if [ "$1" = "install" ]; then
  target="$3"
  version="${target#playwright@}"
  mkdir -p node_modules/playwright
  printf '{"version":"%s"}\n' "$version" > node_modules/playwright/package.json
elif [ "$1" = "exec" ]; then
  : > chromium-ready
fi
`)
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

	if _, err := EnsurePlaywrightVersion(base, "1.61.0", nil, nil); err != nil {
		t.Fatal(err)
	}
	if got := readTestFile(t, logPath); !strings.Contains(got, "npm install --save-exact playwright@1.61.0") {
		t.Fatalf("mismatched package was not reinstalled:\n%s", got)
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
