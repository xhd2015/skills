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

// SkillHeader is standard SKILL.md frontmatter.
//
// Deprecated: use skillcmd.SkillHeader.
type SkillHeader = skillcmd.SkillHeader

// ErrSkillVersionMissing reports that metadata.version is absent or empty.
//
// Deprecated: use skillcmd.ErrSkillVersionMissing.
var ErrSkillVersionMissing = skillcmd.ErrSkillVersionMissing

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

// ParseSkillHeader parses standard SKILL.md YAML frontmatter.
//
// Deprecated: use skillcmd.ParseSkillHeader.
func ParseSkillHeader(content string) (SkillHeader, error) {
	return skillcmd.ParseSkillHeader(content)
}

// SkillVersion returns metadata.version from standard SKILL.md frontmatter.
//
// Deprecated: use skillcmd.SkillVersion.
func SkillVersion(content string) (string, error) {
	return skillcmd.SkillVersion(content)
}

// FormatHeaderWithDelimiters returns the YAML frontmatter block including delimiter lines.
//
// Deprecated: use skillcmd.FormatHeaderWithDelimiters.
func FormatHeaderWithDelimiters(content string) (string, error) {
	return skillcmd.FormatHeaderWithDelimiters(content)
}
