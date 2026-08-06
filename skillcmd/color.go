package skillcmd

import (
	"os"

	"golang.org/x/term"
)

// ColorMode selects how ANSI color is resolved for update output.
type ColorMode int

const (
	// ColorAuto enables color only when stdout is a TTY and NO_COLOR is unset/empty.
	ColorAuto ColorMode = iota
	// ColorAlways forces ANSI on (ignores TTY and NO_COLOR).
	ColorAlways
	// ColorNever forces ANSI off.
	ColorNever
)

const (
	ansiReset = "\x1b[0m"
	ansiRed   = "\x1b[31m"
	ansiGreen = "\x1b[32m"
	ansiYellow = "\x1b[33m"
	ansiGray  = "\x1b[90m"
)

// ResolveColor returns whether ANSI escapes should be emitted.
// stdoutIsTTY should come from term.IsTerminal on the real stdout fd.
// noColorEnv is os.Getenv("NO_COLOR") (empty string if unset).
func ResolveColor(mode ColorMode, stdoutIsTTY bool, noColorEnv string) bool {
	switch mode {
	case ColorAlways:
		return true
	case ColorNever:
		return false
	default: // ColorAuto
		if noColorEnv != "" {
			return false
		}
		return stdoutIsTTY
	}
}

// stdoutIsTTY reports whether os.Stdout is a terminal.
func stdoutIsTTY() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}

type colorStyle struct{ enabled bool }

func newColorStyle(mode ColorMode) colorStyle {
	return colorStyle{enabled: ResolveColor(mode, stdoutIsTTY(), os.Getenv("NO_COLOR"))}
}

func (c colorStyle) wrap(code, s string) string {
	if !c.enabled || s == "" {
		return s
	}
	return code + s + ansiReset
}

func (c colorStyle) green(s string) string  { return c.wrap(ansiGreen, s) }
func (c colorStyle) red(s string) string    { return c.wrap(ansiRed, s) }
func (c colorStyle) yellow(s string) string { return c.wrap(ansiYellow, s) }
func (c colorStyle) gray(s string) string   { return c.wrap(ansiGray, s) }

// colorStatus paints status tokens: updated/would update green; up to date/not installed gray.
func (c colorStyle) colorStatus(status string) string {
	switch status {
	case "updated", "would update":
		return c.green(status)
	case "up to date", "not installed":
		return c.gray(status)
	default:
		return status
	}
}
