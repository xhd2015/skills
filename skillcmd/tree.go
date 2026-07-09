package skillcmd

import (
	"fmt"
	"io/fs"
	"path"
	"sort"
	"strings"
)

// validatePathSegments rejects empty, ".", and ".." path segments.
func validatePathSegments(segments []string) error {
	for _, s := range segments {
		if s == "" || s == "." || s == ".." {
			return fmt.Errorf("invalid path segment: %q", s)
		}
	}
	return nil
}

// loadTreeSkill reads path/SKILL.md from treeFS. topicPath is slash-separated
// (e.g. "a/b"). Empty path is not valid here — callers use RootContent for root.
func loadTreeSkill(treeFS fs.FS, topicPath string) (string, error) {
	topicPath = strings.Trim(topicPath, "/")
	if topicPath == "" {
		return "", fmt.Errorf("empty topic path")
	}
	segments := strings.Split(topicPath, "/")
	if err := validatePathSegments(segments); err != nil {
		return "", err
	}
	embedPath := path.Join(topicPath, "SKILL.md")
	data, err := fs.ReadFile(treeFS, embedPath)
	if err != nil {
		return "", fmt.Errorf("read skill %s: %w", topicPath, err)
	}
	return string(data), nil
}

// ListTreeTopics returns sorted slash-separated topic paths for every nested
// path/SKILL.md under treeFS (e.g. "flags-parsing", "flags-parsing/types").
// Root-level SKILL.md is not a topic path.
func ListTreeTopics(treeFS fs.FS) ([]string, error) {
	if treeFS == nil {
		return nil, nil
	}
	var topics []string
	err := fs.WalkDir(treeFS, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		p = path.Clean(p)
		if path.Base(p) != "SKILL.md" {
			return nil
		}
		if p == "SKILL.md" || p == "./SKILL.md" {
			return nil
		}
		rel := strings.TrimSuffix(p, "/SKILL.md")
		rel = strings.TrimSuffix(rel, "SKILL.md")
		rel = strings.Trim(rel, "/")
		if rel == "" {
			return nil
		}
		topics = append(topics, rel)
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(topics)
	return topics, nil
}

// FormatTopicIndex prints a hierarchical topic list (indented by path depth).
func FormatTopicIndex(topics []string) string {
	if len(topics) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString("Available topics:\n")
	for _, t := range topics {
		depth := strings.Count(t, "/")
		indent := strings.Repeat("  ", depth)
		label := t
		if idx := strings.LastIndex(t, "/"); idx >= 0 {
			label = t[idx+1:]
		}
		b.WriteString("  ")
		b.WriteString(indent)
		b.WriteString("- ")
		b.WriteString(label)
		b.WriteByte('\n')
	}
	return b.String()
}

// collectTreeSkillFiles walks treeFS and returns every nested path/SKILL.md
// (excluding a root-level SKILL.md) as InstallFile entries.
func collectTreeSkillFiles(treeFS fs.FS) ([]InstallFile, error) {
	var files []InstallFile
	err := fs.WalkDir(treeFS, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		// normalize for path.Base / comparisons (WalkDir uses /)
		p = path.Clean(p)
		if path.Base(p) != "SKILL.md" {
			return nil
		}
		if p == "SKILL.md" || p == "./SKILL.md" {
			return nil
		}
		data, err := fs.ReadFile(treeFS, p)
		if err != nil {
			return err
		}
		files = append(files, InstallFile{
			Path:    p,
			Content: data,
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}
