package skillcmd

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// Entry is a single YAML frontmatter field.
type Entry struct {
	Name  string
	Value string
}

// Entries is an ordered list of YAML frontmatter fields.
type Entries []Entry

// Get returns the value for a header field by exact key name, or "" if absent.
func (e Entries) Get(name string) string {
	for _, entry := range e {
		if entry.Name == name {
			return entry.Value
		}
	}
	return ""
}

// TrimHeader strips the YAML frontmatter header from a skill file.
// The header is delimited by "---\n" at the start and a closing line of
// optional whitespace plus "---". Returns the body without the header,
// trimmed of leading/trailing whitespace.
func TrimHeader(content string) string {
	s := content
	if strings.HasPrefix(s, "---\n") {
		rest := s[4:]
		if end, after := findClosingFrontmatter(rest); end >= 0 {
			s = rest[after:]
			if strings.HasPrefix(s, "\n") {
				s = s[1:]
			}
		}
	}
	return strings.TrimSpace(s)
}

// GetHeader returns the inner YAML text between --- delimiters without the delimiters.
// Closing --- may be indented (e.g. tabs from generated test sources).
func GetHeader(content string) (string, error) {
	if !strings.HasPrefix(content, "---\n") {
		return "", fmt.Errorf("missing YAML frontmatter header")
	}
	rest := content[4:]
	end, _ := findClosingFrontmatter(rest)
	if end < 0 {
		return "", fmt.Errorf("missing closing YAML frontmatter delimiter")
	}
	return strings.TrimSpace(rest[:end]), nil
}

// findClosingFrontmatter locates a line that is only optional whitespace plus ---.
// It returns (index of content before that line's preceding newline, index after
// the --- token). end is -1 when not found.
//
// For rest "name: x\n---\nbody", end=6 (at '\n' before ---? actually end is
// length of header text "name: x"), after points past "---".
func findClosingFrontmatter(rest string) (end, after int) {
	// Fast path: classic unindented closer.
	if idx := strings.Index(rest, "\n---"); idx >= 0 {
		// Ensure the line is only --- (optional trailing spaces/CR), not ---something.
		line := rest[idx+1:]
		if nl := strings.IndexByte(line, '\n'); nl >= 0 {
			line = line[:nl]
		}
		line = strings.TrimSuffix(line, "\r")
		if strings.TrimSpace(line) == "---" {
			return idx, idx + 1 + len(line)
		}
	}
	// Line-scan: allow leading whitespace on the closing --- line.
	offset := 0
	for offset <= len(rest) {
		nl := strings.IndexByte(rest[offset:], '\n')
		var line string
		lineStart := offset
		next := -1
		if nl < 0 {
			line = rest[lineStart:]
		} else {
			line = rest[lineStart : lineStart+nl]
			next = lineStart + nl + 1
		}
		trimmed := strings.TrimSuffix(line, "\r")
		if strings.TrimSpace(trimmed) == "---" {
			// Header text ends at the character before this line (the prior newline),
			// or at 0 if --- is the first line of rest.
			if lineStart == 0 {
				return 0, len(line)
			}
			// lineStart points just after the newline that ends the previous line;
			// header excludes that newline → end = lineStart-1 when we want
			// rest[:end] without the newline. Prefer end = lineStart-1 only if
			// rest[lineStart-1]=='\n'; use lineStart-1 to match classic Index.
			end = lineStart - 1
			if end < 0 {
				end = 0
			}
			return end, lineStart + len(line)
		}
		if next < 0 {
			return -1, -1
		}
		offset = next
	}
	return -1, -1
}

// ParseHeader parses YAML header text into ordered entries.
func ParseHeader(header string) (Entries, error) {
	var root yaml.Node
	if err := yaml.Unmarshal([]byte(header), &root); err != nil {
		return nil, err
	}
	if root.Kind != yaml.DocumentNode || len(root.Content) == 0 {
		return nil, fmt.Errorf("empty header")
	}
	mapNode := root.Content[0]
	if mapNode.Kind != yaml.MappingNode {
		return nil, fmt.Errorf("header must be a mapping")
	}
	entries := make(Entries, 0, len(mapNode.Content)/2)
	for i := 0; i < len(mapNode.Content); i += 2 {
		keyNode := mapNode.Content[i]
		valNode := mapNode.Content[i+1]
		var value string
		if err := valNode.Decode(&value); err != nil {
			return nil, err
		}
		entries = append(entries, Entry{
			Name:  keyNode.Value,
			Value: strings.TrimSpace(value),
		})
	}
	return entries, nil
}

// FormatHeaderWithDelimiters returns the YAML frontmatter block including delimiter lines.
func FormatHeaderWithDelimiters(content string) (string, error) {
	header, err := GetHeader(content)
	if err != nil {
		return "", err
	}
	return "---\n" + header + "\n---\n", nil
}