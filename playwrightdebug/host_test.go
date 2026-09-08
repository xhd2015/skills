package playwrightdebug

import (
	"os"
	"strings"
	"testing"
)

func TestParseOSRelease(t *testing.T) {
	id, ver := parseOSRelease(strings.NewReader(`PRETTY_NAME="Debian GNU/Linux 11 (bullseye)"
NAME="Debian GNU/Linux"
VERSION_ID="11"
ID=debian
VERSION_CODENAME=bullseye
`))
	if id != "debian" || ver != "11" {
		t.Fatalf("parseOSRelease = %q %q, want debian 11", id, ver)
	}

	id, ver = parseOSRelease(strings.NewReader("ID=ubuntu\nVERSION_ID=22.04\n"))
	if id != "ubuntu" || ver != "22.04" {
		t.Fatalf("parseOSRelease ubuntu = %q %q", id, ver)
	}
}

func TestHostOSIsDebian11Amd64(t *testing.T) {
	if !(HostOS{ID: "debian", VersionID: "11", Arch: "amd64"}).IsDebian11Amd64() {
		t.Fatal("amd64 debian 11 should match")
	}
	if !(HostOS{ID: "Debian", VersionID: "11", Arch: "x86_64"}).IsDebian11Amd64() {
		t.Fatal("x86_64 debian 11 should match")
	}
	if (HostOS{ID: "debian", VersionID: "12", Arch: "amd64"}).IsDebian11Amd64() {
		t.Fatal("debian 12 should not match")
	}
	if (HostOS{ID: "debian", VersionID: "11", Arch: "arm64"}).IsDebian11Amd64() {
		t.Fatal("arm64 should not match")
	}
}

func TestBuildRunPlanAutoPinsDebian11CacheDir(t *testing.T) {
	prev := lookupHostOS
	lookupHostOS = func() HostOS {
		return HostOS{ID: "debian", VersionID: "11", Arch: "amd64"}
	}
	t.Cleanup(func() { lookupHostOS = prev })

	script := t.TempDir() + "/probe.js"
	if err := os.WriteFile(script, []byte("console.log('ok')"), 0o644); err != nil {
		t.Fatal(err)
	}
	base := t.TempDir()
	var stderr strings.Builder
	plan, err := BuildRunPlan(RunOptions{
		ScriptPath: script,
		CacheDir:   base,
		SkipEnsure: true,
		Stderr:     &stderr,
	})
	if err != nil {
		t.Fatal(err)
	}
	if plan.PlaywrightVersion != LastDebian11PlaywrightVersion {
		t.Fatalf("plan version = %q", plan.PlaywrightVersion)
	}
	wantCache, err := VersionedCacheDir(base, LastDebian11PlaywrightVersion)
	if err != nil {
		t.Fatal(err)
	}
	if plan.CacheDir != wantCache {
		t.Fatalf("cache = %q, want %q", plan.CacheDir, wantCache)
	}
	if !strings.Contains(stderr.String(), "warning: debian 11 x64") {
		t.Fatalf("missing warning: %q", stderr.String())
	}
}

func TestResolvePlaywrightVersion(t *testing.T) {
	debian11 := HostOS{ID: "debian", VersionID: "11", Arch: "amd64"}
	got, auto := ResolvePlaywrightVersion("", debian11)
	if got != LastDebian11PlaywrightVersion || !auto {
		t.Fatalf("empty on debian11 = %q auto=%v, want %s true", got, auto, LastDebian11PlaywrightVersion)
	}
	got, auto = ResolvePlaywrightVersion("1.55.0", debian11)
	if got != "1.55.0" || auto {
		t.Fatalf("explicit on debian11 = %q auto=%v", got, auto)
	}
	got, auto = ResolvePlaywrightVersion("", HostOS{ID: "ubuntu", VersionID: "22.04", Arch: "amd64"})
	if got != "" || auto {
		t.Fatalf("empty on ubuntu = %q auto=%v, want unpinned", got, auto)
	}
}
