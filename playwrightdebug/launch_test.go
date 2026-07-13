package playwrightdebug

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExtractLaunchFlags_ExtensionAndRest(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "manifest.json"), []byte(`{"manifest_version":3}`), 0o644); err != nil {
		t.Fatal(err)
	}

	opts, rest, err := ExtractLaunchFlags([]string{
		"--extension", dir,
		"run", "script.js",
		"--foo", "bar",
	})
	if err != nil {
		t.Fatalf("ExtractLaunchFlags: %v", err)
	}
	if !opts.HasExtension() {
		t.Fatal("expected extension")
	}
	if opts.ExtensionPaths[0] != dir && filepath.Clean(opts.ExtensionPaths[0]) != filepath.Clean(dir) {
		// Abs may expand; ensure suffix match
		if filepath.Base(opts.ExtensionPaths[0]) != filepath.Base(dir) {
			t.Fatalf("ext path=%q want under %q", opts.ExtensionPaths[0], dir)
		}
	}
	if opts.Headed == nil || !*opts.Headed {
		// default headed for extension is applied in ApplyEnv, not Extract
	}
	wantRest := []string{"run", "script.js", "--foo", "bar"}
	if len(rest) != len(wantRest) {
		t.Fatalf("rest=%v want %v", rest, wantRest)
	}
	for i := range wantRest {
		if rest[i] != wantRest[i] {
			t.Fatalf("rest[%d]=%q want %q", i, rest[i], wantRest[i])
		}
	}
}

func TestExtractLaunchFlags_MissingManifest(t *testing.T) {
	dir := t.TempDir()
	_, _, err := ExtractLaunchFlags([]string{"--extension", dir})
	if err == nil {
		t.Fatal("expected error for missing manifest.json")
	}
}

func TestApplyEnv_ExtensionDefaultsHeaded(t *testing.T) {
	dir := t.TempDir()
	_ = os.WriteFile(filepath.Join(dir, "manifest.json"), []byte(`{}`), 0o644)
	opts, err := LaunchOptions{ExtensionPaths: []string{dir}}.Normalize()
	if err != nil {
		t.Fatal(err)
	}
	env := opts.ApplyEnv(nil)
	m := map[string]string{}
	for _, e := range env {
		if i := indexByte(e, '='); i >= 0 {
			m[e[:i]] = e[i+1:]
		}
	}
	if m[EnvLaunchMode] != "extension" {
		t.Fatalf("mode=%q", m[EnvLaunchMode])
	}
	if m[EnvHeaded] != "1" {
		t.Fatalf("headed=%q want 1", m[EnvHeaded])
	}
	if m[EnvExtensionPaths] == "" {
		t.Fatal("missing extension paths")
	}
}

func TestExtractLaunchFlags_UserDataDirOnly(t *testing.T) {
	profile := t.TempDir()
	opts, rest, err := ExtractLaunchFlags([]string{
		"--user-data-dir", profile,
		"-e", "console.log(1)",
	})
	if err != nil {
		t.Fatalf("ExtractLaunchFlags: %v", err)
	}
	if opts.UserDataDir == "" {
		t.Fatal("expected UserDataDir")
	}
	if filepath.Clean(opts.UserDataDir) != filepath.Clean(profile) {
		// Abs may expand
		if filepath.Base(opts.UserDataDir) != filepath.Base(profile) {
			t.Fatalf("UserDataDir=%q want under %q", opts.UserDataDir, profile)
		}
	}
	if opts.HasExtension() {
		t.Fatal("did not expect extension")
	}
	wantRest := []string{"-e", "console.log(1)"}
	if len(rest) != len(wantRest) {
		t.Fatalf("rest=%v want %v", rest, wantRest)
	}
}

func TestApplyEnv_UserDataDirWithoutExtension(t *testing.T) {
	profile := t.TempDir()
	opts, err := LaunchOptions{UserDataDir: profile}.Normalize()
	if err != nil {
		t.Fatal(err)
	}
	env := opts.ApplyEnv(nil)
	m := map[string]string{}
	for _, e := range env {
		if i := indexByte(e, '='); i >= 0 {
			m[e[:i]] = e[i+1:]
		}
	}
	if m[EnvLaunchMode] != "default" {
		t.Fatalf("mode=%q want default (profile alone is not extension mode)", m[EnvLaunchMode])
	}
	if m[EnvUserDataDir] == "" {
		t.Fatal("missing user data dir env")
	}
	if m[EnvHeaded] != "0" {
		t.Fatalf("headed=%q want 0 (default stays headless unless --headed)", m[EnvHeaded])
	}
}

func indexByte(s string, c byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			return i
		}
	}
	return -1
}
