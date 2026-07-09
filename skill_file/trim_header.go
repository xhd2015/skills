// Package skill_file provides YAML frontmatter helpers for SKILL.md files.
//
// Deprecated: use package github.com/xhd2015/skills/skillcmd instead.
// This package is a thin re-export shim and will be removed in a future version.
package skill_file

import "github.com/xhd2015/skills/skillcmd"

// Entry is a single YAML frontmatter field.
//
// Deprecated: use skillcmd.Entry.
type Entry = skillcmd.Entry

// Entries is an ordered list of YAML frontmatter fields.
//
// Deprecated: use skillcmd.Entries.
type Entries = skillcmd.Entries

// TrimHeader strips the YAML frontmatter header from a skill file.
//
// Deprecated: use skillcmd.TrimHeader.
func TrimHeader(content string) string {
	return skillcmd.TrimHeader(content)
}

// GetHeader returns the inner YAML text between --- delimiters without the delimiters.
//
// Deprecated: use skillcmd.GetHeader.
func GetHeader(content string) (string, error) {
	return skillcmd.GetHeader(content)
}

// ParseHeader parses YAML header text into ordered entries.
//
// Deprecated: use skillcmd.ParseHeader.
func ParseHeader(header string) (Entries, error) {
	return skillcmd.ParseHeader(header)
}

// FormatHeaderWithDelimiters returns the YAML frontmatter block including delimiter lines.
//
// Deprecated: use skillcmd.FormatHeaderWithDelimiters.
func FormatHeaderWithDelimiters(content string) (string, error) {
	return skillcmd.FormatHeaderWithDelimiters(content)
}
