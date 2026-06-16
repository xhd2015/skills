package skill_file

import "strings"

// TrimHeader strips the YAML frontmatter header from a skill file.
// The header is delimited by "---\n" at the start and "\n---\n" to close.
// Returns the body without the header, trimmed of leading/trailing whitespace.
func TrimHeader(content string) string {
	s := content
	if strings.HasPrefix(s, "---\n") {
		rest := s[4:]
		if idx := strings.Index(rest, "\n---\n"); idx >= 0 {
			s = rest[idx+5:]
			if strings.HasPrefix(s, "\n") {
				s = s[1:]
			}
		}
	}
	return strings.TrimSpace(s)
}
