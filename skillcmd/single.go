package skillcmd

import (
	"fmt"
	"io/fs"
	"strings"
)

// SingleSkill hosts one skill definition with optional nested TreeFS of
// path/TOPIC.md files and optional ExtraFiles for install.
type SingleSkill struct {
	Name        string
	RootContent string
	TreeFS      fs.FS // optional; path "a/b" → "a/b/TOPIC.md"
	ExtraFiles  []InstallFile
	Usage       string
	// Help is printed for skill -h/--help (and --show|--list|--version with --help).
	// When empty, DefaultSingleSkillHelp(Usage, Name) is used.
	Help string
}

// DefaultSingleSkillHelp returns a generic skill subcommand help block.
// For multi-topic skills (TreeFS set), --list and --help also list topic paths.
func DefaultSingleSkillHelp(usage, skillName string) string {
	usage = strings.TrimSpace(usage)
	if usage == "" {
		usage = "skill --install"
	}
	name := strings.TrimSpace(skillName)
	if name == "" {
		name = "<skill-name>"
	}
	return fmt.Sprintf(`Usage: skill --show [--header] [<topic-path>]
       skill <topic-path> --show [--header]
       skill --version
       skill --install [OPTIONS] [<dir>]
       skill <topic-path> --install [OPTIONS] [<dir>]
       skill --list

Show the embedded SKILL.md (root) or a nested topic path (path/TOPIC.md).
Install copies SKILL.md (and nested TOPIC.md topics) to agent skill directories.
A leading <topic-path> before --install is ignored for destination (whole skill
is installed); use --dir or <dir> for the target.
List prints the skill name (%s); when nested topics exist, also lists all topic paths.

Install usage: %s [OPTIONS] [<dir>]
  Run skill --install --help for install flags (--global, --cursor, --dir, …).

Options:
  --show [--header] [path]   Print skill or topic content
  --version                  Print root SKILL.md metadata.version
  --install [OPTIONS] [dir]  Install skill files (see --install --help)
  --list                     Print skill name and available topics (if any)
  -h, --help                 Show this help message
`, name, usage)
}

// Handle runs list/show/version/install for this skill using ParseSkillArgs.
func (s *SingleSkill) Handle(args []string) error {
	parsed, err := ParseSkillArgs(args)
	if err != nil {
		return err
	}
	switch parsed.Action {
	case ActionHelp:
		fmt.Print(s.helpText())
		return nil
	case ActionList:
		return s.handleList()
	case ActionShow:
		return s.handleShow(parsed.Header, parsed.Rest)
	case ActionVersion:
		return s.handleVersion(parsed.Rest)
	case ActionInstall:
		return s.handleInstall(parsed.Rest)
	default:
		return fmt.Errorf("unknown action: %s", parsed.Action)
	}
}

func (s *SingleSkill) helpText() string {
	base := strings.TrimSpace(s.Help)
	if base == "" {
		base = DefaultSingleSkillHelp(s.Usage, s.Name)
	}
	if !strings.HasSuffix(base, "\n") {
		base += "\n"
	}
	topics, err := s.topicPaths()
	if err != nil || len(topics) == 0 {
		return base
	}
	return base + "\n" + FormatTopicIndex(topics)
}

// handleList prints the skill name; for multi-topic (TreeFS) skills also lists
// every nested topic path (full slash paths, one per line after the name).
func (s *SingleSkill) handleList() error {
	fmt.Println(s.Name)
	topics, err := s.topicPaths()
	if err != nil {
		return err
	}
	for _, t := range topics {
		fmt.Println(t)
	}
	return nil
}

func (s *SingleSkill) topicPaths() ([]string, error) {
	if s.TreeFS == nil {
		return nil, nil
	}
	return ListTreeTopics(s.TreeFS)
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

func (s *SingleSkill) handleVersion(rest []string) error {
	if len(rest) > 0 {
		return fmt.Errorf("unexpected arguments: %v", rest)
	}
	return printSkillVersion(s.Name, s.RootContent)
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
	rest, err := s.peelLeadingTopicForInstall(rest)
	if err != nil {
		return err
	}
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

// peelLeadingTopicForInstall drops a known TreeFS topic from Rest when the
// caller used show-style order (`skill <topic> --install …`) and also supplied
// an install destination (--dir or another positional). Install always copies
// the whole skill; the topic token is not a destination.
//
// If the only non-flag token is a topic name and --dir is absent, it is left
// alone so `skill --install <dir>` can still use a directory whose basename
// happens to match a topic.
func (s *SingleSkill) peelLeadingTopicForInstall(rest []string) ([]string, error) {
	if s.TreeFS == nil || len(rest) == 0 {
		return rest, nil
	}
	topics, err := ListTreeTopics(s.TreeFS)
	if err != nil {
		return nil, err
	}
	known := make(map[string]struct{}, len(topics))
	for _, t := range topics {
		known[t] = struct{}{}
	}

	topicIdx := -1
	hasDirFlag := false
	otherPositionals := 0
	for i := 0; i < len(rest); i++ {
		arg := rest[i]
		if arg == "--dir" {
			hasDirFlag = true
			// skip flag value if present
			if i+1 < len(rest) && !strings.HasPrefix(rest[i+1], "-") {
				i++
			}
			continue
		}
		if strings.HasPrefix(arg, "-") {
			continue
		}
		cand := strings.Trim(arg, "/")
		if topicIdx < 0 {
			if _, ok := known[cand]; ok {
				topicIdx = i
				continue
			}
		}
		otherPositionals++
	}
	if topicIdx < 0 || (!hasDirFlag && otherPositionals == 0) {
		return rest, nil
	}
	out := make([]string, 0, len(rest)-1)
	out = append(out, rest[:topicIdx]...)
	out = append(out, rest[topicIdx+1:]...)
	return out, nil
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
