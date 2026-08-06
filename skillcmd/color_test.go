package skillcmd

import "testing"

func TestResolveColor(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name       string
		mode       ColorMode
		tty        bool
		noColorEnv string
		want       bool
	}{
		{"always non-tty", ColorAlways, false, "", true},
		{"always with NO_COLOR", ColorAlways, true, "1", true},
		{"never tty", ColorNever, true, "", false},
		{"auto tty", ColorAuto, true, "", true},
		{"auto pipe", ColorAuto, false, "", false},
		{"auto NO_COLOR", ColorAuto, true, "1", false},
		{"auto empty NO_COLOR", ColorAuto, true, "", true},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := ResolveColor(tc.mode, tc.tty, tc.noColorEnv)
			if got != tc.want {
				t.Fatalf("ResolveColor(%v, tty=%v, NO_COLOR=%q)=%v want %v",
					tc.mode, tc.tty, tc.noColorEnv, got, tc.want)
			}
		})
	}
}

func TestColorStyleWrapDisabled(t *testing.T) {
	t.Parallel()
	c := colorStyle{enabled: false}
	if got := c.green("updated"); got != "updated" {
		t.Fatalf("disabled green: %q", got)
	}
}

func TestColorStyleWrapEnabled(t *testing.T) {
	t.Parallel()
	c := colorStyle{enabled: true}
	got := c.green("updated")
	if got == "updated" || !containsESC(got) {
		t.Fatalf("enabled green missing ANSI: %q", got)
	}
	if got != ansiGreen+"updated"+ansiReset {
		t.Fatalf("unexpected wrap: %q", got)
	}
}

func containsESC(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == 0x1b {
			return true
		}
	}
	return false
}
