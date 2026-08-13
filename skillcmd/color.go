package skillcmd

import "github.com/xhd2015/dot-pkgs/go-pkgs/terminal/color"

// colorStatus paints status tokens: updated/would update green; up to date/not installed gray.
func colorStatus(s color.Style, status string) string {
	switch status {
	case "updated", "would update":
		return s.Green(status)
	case "up to date", "not installed":
		return s.Gray(status)
	default:
		return status
	}
}
