package playwrightdebug

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
)

// LastDebian11PlaywrightVersion is the last Playwright release that still
// shipped Chromium for debian11-x64 (1.62 dropped Debian 11).
const LastDebian11PlaywrightVersion = "1.61.0"

// HostOS is the host identity used to pick a default Playwright version.
type HostOS struct {
	ID        string
	VersionID string
	Arch      string
}

func (h HostOS) IsDebian11Amd64() bool {
	id := strings.ToLower(strings.TrimSpace(h.ID))
	ver := strings.TrimSpace(h.VersionID)
	arch := strings.ToLower(strings.TrimSpace(h.Arch))
	if id != "debian" || ver != "11" {
		return false
	}
	return arch == "amd64" || arch == "x86_64"
}

var lookupHostOS = detectHostOS

func detectHostOS() HostOS {
	return HostOS{
		Arch: runtime.GOARCH,
	}.withRelease(readOSReleaseFile("/etc/os-release"))
}

func (h HostOS) withRelease(id, versionID string) HostOS {
	h.ID = id
	h.VersionID = versionID
	return h
}

func readOSReleaseFile(path string) (id, versionID string) {
	f, err := os.Open(path)
	if err != nil {
		return "", ""
	}
	defer f.Close()
	return parseOSRelease(f)
}

func parseOSRelease(r io.Reader) (id, versionID string) {
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, val, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		val = strings.Trim(val, `"'`)
		switch key {
		case "ID":
			id = val
		case "VERSION_ID":
			versionID = val
		}
	}
	return id, versionID
}

// ResolvePlaywrightVersion returns the Playwright package version to install.
// A non-empty requested version always wins. When requested is empty and the
// host is debian11-x64, it pins LastDebian11PlaywrightVersion.
func ResolvePlaywrightVersion(requested string, host HostOS) (resolved string, autoPinned bool) {
	if strings.TrimSpace(requested) != "" {
		return requested, false
	}
	if host.IsDebian11Amd64() {
		return LastDebian11PlaywrightVersion, true
	}
	return "", false
}

func applyResolvedPlaywrightVersion(requested string, stderr io.Writer) string {
	resolved, auto := ResolvePlaywrightVersion(requested, lookupHostOS())
	if auto {
		if stderr == nil {
			stderr = os.Stderr
		}
		fmt.Fprintf(stderr, "warning: debian 11 x64: using Playwright %s (last version that supports this OS)\n", resolved)
	}
	return resolved
}
