package skill_file

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

// GetHeader returns the inner YAML text between --- delimiters without the delimiters.
func GetHeader(content string) (string, error) {
	if !strings.HasPrefix(content, "---\n") {
		return "", fmt.Errorf("missing YAML frontmatter header")
	}
	rest := content[4:]
	idx := strings.Index(rest, "\n---")
	if idx < 0 {
		return "", fmt.Errorf("missing closing YAML frontmatter delimiter")
	}
	return strings.TrimSpace(rest[:idx]), nil
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