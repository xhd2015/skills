package skillcmd

import (
	"errors"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// SkillHeader is standard SKILL.md frontmatter. Metadata values are strings as
// required by the Agent Skills specification.
type SkillHeader struct {
	Name          string
	Description   string
	License       string
	Compatibility string
	AllowedTools  string
	Metadata      map[string]string
}

// ErrSkillVersionMissing reports that metadata.version is absent or empty.
var ErrSkillVersionMissing = errors.New("metadata.version is not declared")

// ParseSkillHeader parses standard SKILL.md YAML frontmatter. Unknown top-level
// fields are ignored; unknown metadata keys are preserved.
func ParseSkillHeader(content string) (SkillHeader, error) {
	headerText, err := GetHeader(content)
	if err != nil {
		return SkillHeader{}, err
	}
	var root yaml.Node
	if err := yaml.Unmarshal([]byte(headerText), &root); err != nil {
		return SkillHeader{}, err
	}
	if root.Kind != yaml.DocumentNode || len(root.Content) == 0 {
		return SkillHeader{}, fmt.Errorf("empty header")
	}
	mapping := root.Content[0]
	if mapping.Kind != yaml.MappingNode {
		return SkillHeader{}, fmt.Errorf("header must be a mapping")
	}

	out := SkillHeader{Metadata: make(map[string]string)}
	for i := 0; i < len(mapping.Content); i += 2 {
		key, value := mapping.Content[i], mapping.Content[i+1]
		switch key.Value {
		case "name":
			if err := value.Decode(&out.Name); err != nil {
				return SkillHeader{}, fmt.Errorf("name: %w", err)
			}
		case "description":
			if err := value.Decode(&out.Description); err != nil {
				return SkillHeader{}, fmt.Errorf("description: %w", err)
			}
		case "license":
			if err := value.Decode(&out.License); err != nil {
				return SkillHeader{}, fmt.Errorf("license: %w", err)
			}
		case "compatibility":
			if err := value.Decode(&out.Compatibility); err != nil {
				return SkillHeader{}, fmt.Errorf("compatibility: %w", err)
			}
		case "allowed-tools":
			if err := value.Decode(&out.AllowedTools); err != nil {
				return SkillHeader{}, fmt.Errorf("allowed-tools: %w", err)
			}
		case "metadata":
			if value.Kind != yaml.MappingNode {
				return SkillHeader{}, fmt.Errorf("metadata must be a mapping from string keys to string values")
			}
			for j := 0; j < len(value.Content); j += 2 {
				metadataKey, metadataValue := value.Content[j], value.Content[j+1]
				if metadataKey.Kind != yaml.ScalarNode || metadataKey.Tag != "!!str" {
					return SkillHeader{}, fmt.Errorf("metadata keys must be strings")
				}
				if metadataValue.Kind != yaml.ScalarNode || metadataValue.Tag != "!!str" {
					return SkillHeader{}, fmt.Errorf("metadata.%s must be a string", metadataKey.Value)
				}
				out.Metadata[metadataKey.Value] = strings.TrimSpace(metadataValue.Value)
			}
		}
	}
	return out, nil
}

// SkillVersion returns metadata.version from standard SKILL.md frontmatter.
func SkillVersion(content string) (string, error) {
	header, err := ParseSkillHeader(content)
	if err != nil {
		return "", err
	}
	version := strings.TrimSpace(header.Metadata["version"])
	if version == "" {
		return "", ErrSkillVersionMissing
	}
	return version, nil
}

func printSkillVersion(name, content string) error {
	version, err := SkillVersion(content)
	if errors.Is(err, ErrSkillVersionMissing) {
		return fmt.Errorf("skill %s has no metadata.version", name)
	}
	if err != nil {
		return fmt.Errorf("skill %s: %w", name, err)
	}
	fmt.Println(version)
	return nil
}
