package skillcmd

import (
	"fmt"
	"io/fs"
	"strings"
)

// SingleSkill hosts one skill definition with optional nested TreeFS of
// path/SKILL.md files and optional ExtraFiles for install.
type SingleSkill struct {
	Name        string
	RootContent string
	TreeFS      fs.FS // optional; path "a/b" → "a/b/SKILL.md"
	ExtraFiles  []InstallFile
	Usage       string
}

// Handle runs list/show/install for this skill using ParseSkillArgs flag surface.
func (s *SingleSkill) Handle(args []string) error {
	parsed, err := ParseSkillArgs(args)
	if err != nil {
		return err
	}
	switch parsed.Action {
	case ActionList:
		fmt.Println(s.Name)
		return nil
	case ActionShow:
		return s.handleShow(parsed.Header, parsed.Rest)
	case ActionInstall:
		return s.handleInstall(parsed.Rest)
	default:
		return fmt.Errorf("unknown action: %s", parsed.Action)
	}
}

func (s *SingleSkill) handleShow(header bool, rest []string) error {
	content, err := s.loadContent(rest)
	if err != nil {
		return err
	}
	if header {
		out, err := FormatHeaderWithDelimiters(content)
		if err != nil {
			return err
		}
		fmt.Print(out)
		return nil
	}
	fmt.Print(content)
	return nil
}

func (s *SingleSkill) loadContent(rest []string) (string, error) {
	if len(rest) == 0 {
		return s.RootContent, nil
	}
	topicPath := strings.Trim(rest[0], "/")
	if len(rest) > 1 {
		return "", fmt.Errorf("unexpected arguments: %v", rest[1:])
	}
	// Validate even when TreeFS is nil so ../ is always rejected.
	segments := strings.Split(topicPath, "/")
	if err := validatePathSegments(segments); err != nil {
		return "", err
	}
	if s.TreeFS == nil {
		return "", fmt.Errorf("unknown topic path: %s", topicPath)
	}
	return loadTreeSkill(s.TreeFS, topicPath)
}

func (s *SingleSkill) handleInstall(rest []string) error {
	extras, err := s.extraFilesForInstall()
	if err != nil {
		return err
	}
	usage := strings.TrimSpace(s.Usage)
	if usage == "" {
		usage = "skill --install"
	}
	return HandleInstall(InstallOptions{
		SkillDirName: s.Name,
		SkillContent: s.RootContent,
		ExtraFiles:   extras,
		Usage:        usage,
	}, rest)
}

func (s *SingleSkill) extraFilesForInstall() ([]InstallFile, error) {
	if s.ExtraFiles != nil {
		return s.ExtraFiles, nil
	}
	if s.TreeFS == nil {
		return nil, nil
	}
	return collectTreeSkillFiles(s.TreeFS)
}
