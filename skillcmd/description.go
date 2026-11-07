package skillcmd

import (
	"strings"

	"gopkg.in/yaml.v3"
)

// ExtractDescription returns the frontmatter description folded into a single
// line: YAML block scalars (>-, |), quoted values, and internal whitespace are
// normalized. Returns "" when the content has no parseable description.
func ExtractDescription(content string) string {
	header, err := GetHeader(content)
	if err != nil {
		return ""
	}
	var meta struct {
		Description string `yaml:"description"`
	}
	if err := yaml.Unmarshal([]byte(header), &meta); err != nil {
		return ""
	}
	return strings.Join(strings.Fields(meta.Description), " ")
}
