package skillcmd

import (
	"fmt"
	"io/fs"
	"strings"
)

// RegisteredSkill is one entry in a multi-skill registry. A skill is a single
// SKILL.md (Content), optionally extended by a TreeFS of path/TOPIC.md topics.
type RegisteredSkill struct {
	Name        string
	Description string // optional; derived from Content frontmatter when empty
	Content     string
	// TreeFS is optional. Topic path "a/b" reads "a/b/TOPIC.md".
	TreeFS     fs.FS
	ExtraFiles []InstallFile
	// Usage is this skill's install help example, e.g.
	// "wrk skill --install wrk-pipeline". When empty, Registry.Usage + name.
	Usage string
}

// Registry is an ordered list of registered skills for multi-skill hosts, so
// one binary can host several independent top-level skills.
//
// A skill reference is "name" or "name/topic-path". Version and install always
// address the skill root; show resolves the topic path against that skill.
type Registry struct {
	Skills []RegisteredSkill
	// Default is the skill used when the reference omits the name (e.g. the
	// host binary's own skill). Empty means every action requires a name.
	Default string
	// Help is printed for skill -h/--help. When empty, DefaultRegistrySkillHelp is used.
	Help string
	// SkillsHelp is printed for skills -h/--help. When empty, DefaultRegistrySkillsHelp is used.
	SkillsHelp string
	// Usage prefixes install help examples (e.g. "agent-pro skill --install").
	Usage string
	// UpdateUsage prefixes skills-update help (default "skills update").
	UpdateUsage string
}

// DefaultRegistrySkillHelp returns help for the multi-skill `skill` entry point.
func DefaultRegistrySkillHelp(usage string) string {
	usage = strings.TrimSpace(usage)
	if usage == "" {
		usage = "skill --install"
	}
	return fmt.Sprintf(`Usage: skill --list
       skill --show [<name>[/<topic>]]
       skill <name>[/<topic>] --show
       skill --version [<name>]
       skill <name> --version
       skill --install [<name>[/<topic>]] [OPTIONS] [<dir>]
       skill <name> --install [OPTIONS] [<dir>]

List registered skills, print a skill's root SKILL.md or a
<topic-path>/TOPIC.md, print metadata.version, or install one whole skill.
The registry default skill is used when <name> is omitted and one is set.

Install example: %s <name> [OPTIONS]
  Run skill --install <name> --help for install flags (--global, --cursor, …).

Options:
  -h, --help    Show this help message
`, usage)
}

// DefaultRegistrySkillsHelp returns help for the multi-skill `skills` entry point.
func DefaultRegistrySkillsHelp() string {
	return `Usage: skills
       skills update [OPTIONS] [<dir>]

With no arguments, list registered skills (same as skill --list).
skills update refreshes already-installed skills only.

  Run skills update --help for update flags (--global, --cursor, …).

Options:
  -h, --help    Show this help message
`
}

// HandleSkill implements skill flag actions: list, show, version, and install.
func (r *Registry) HandleSkill(args []string) error {
	parsed, err := ParseSkillArgs(args)
	if err != nil {
		return err
	}
	switch parsed.Action {
	case ActionHelp:
		fmt.Print(r.skillHelpText())
		return nil
	case ActionList:
		return r.listSkills()
	case ActionShow:
		return r.handleShow(parsed.Header, parsed.Rest)
	case ActionVersion:
		return r.handleVersion(parsed.Rest)
	case ActionInstall:
		return r.handleInstall(parsed.Rest)
	default:
		return fmt.Errorf("unknown action: %s", parsed.Action)
	}
}

// HandleSkills implements the skills entry point: empty → list; update →
// HandleUpdateMany; otherwise delegates to HandleSkill.
func (r *Registry) HandleSkills(args []string) error {
	if len(args) == 0 {
		return r.listSkills()
	}
	if args[0] == "-h" || args[0] == "--help" {
		fmt.Print(r.skillsHelpText())
		return nil
	}
	if args[0] == "update" {
		return r.handleUpdateMany(args[1:])
	}
	return r.HandleSkill(args)
}

func (r *Registry) skillHelpText() string {
	base := r.Help
	if strings.TrimSpace(base) == "" {
		base = DefaultRegistrySkillHelp(r.Usage)
	}
	if !strings.HasSuffix(base, "\n") {
		base += "\n"
	}
	return base + r.topicIndexText()
}

// topicIndexText lists "name/topic" paths for tree skills. Empty when every
// registered skill is a single SKILL.md.
func (r *Registry) topicIndexText() string {
	var lines []string
	for _, sk := range r.Skills {
		if sk.TreeFS == nil {
			continue
		}
		topics, err := ListTreeTopics(sk.TreeFS)
		if err != nil {
			continue
		}
		for _, t := range topics {
			lines = append(lines, sk.Name+"/"+t)
		}
	}
	if len(lines) == 0 {
		return ""
	}
	return "\n" + FormatTopicIndex(lines)
}

func (r *Registry) skillsHelpText() string {
	if strings.TrimSpace(r.SkillsHelp) != "" {
		return r.SkillsHelp
	}
	return DefaultRegistrySkillsHelp()
}

func (r *Registry) listSkills() error {
	for _, sk := range r.Skills {
		desc := strings.TrimSpace(sk.Description)
		if desc == "" {
			desc = ExtractDescription(sk.Content)
		}
		if desc != "" {
			fmt.Printf("%s - %s\n", sk.Name, desc)
		} else {
			fmt.Println(sk.Name)
		}
		if sk.TreeFS == nil {
			continue
		}
		topics, err := ListTreeTopics(sk.TreeFS)
		if err != nil {
			return err
		}
		for _, t := range topics {
			fmt.Printf("%s/%s\n", sk.Name, t)
		}
	}
	return nil
}

func (r *Registry) find(name string) (RegisteredSkill, bool) {
	for _, sk := range r.Skills {
		if sk.Name == name {
			return sk, true
		}
	}
	return RegisteredSkill{}, false
}

func (r *Registry) namesText() string {
	names := make([]string, 0, len(r.Skills))
	for _, sk := range r.Skills {
		names = append(names, sk.Name)
	}
	return strings.Join(names, ", ")
}

// single adapts one registry entry to SingleSkill so topic loading, install
// file collection, and install destination naming stay in one implementation.
func (r *Registry) single(sk RegisteredSkill) *SingleSkill {
	return &SingleSkill{
		Name:        sk.Name,
		RootContent: sk.Content,
		TreeFS:      sk.TreeFS,
		ExtraFiles:  sk.ExtraFiles,
		Usage:       r.installUsage(sk),
	}
}

func (r *Registry) installUsage(sk RegisteredSkill) string {
	if usage := strings.TrimSpace(sk.Usage); usage != "" {
		return usage
	}
	base := strings.TrimSpace(r.Usage)
	if base == "" {
		base = "skill --install"
	}
	return base + " " + sk.Name
}

// resolve maps a positional skill reference ("name" or "name/topic") to a
// registered skill. Only show hands the topic back to the caller; version and
// install address the skill root, so their topic suffix is dropped.
func (r *Registry) resolve(rest []string, action Action) (RegisteredSkill, []string, error) {
	rest = append([]string(nil), rest...)
	name := ""
	if len(rest) > 0 && !strings.HasPrefix(rest[0], "-") {
		name = rest[0]
		rest = rest[1:]
	} else if def := strings.TrimSpace(r.Default); def != "" {
		name = def
	}
	if name == "" {
		return RegisteredSkill{}, nil, fmt.Errorf("expected skill name for --%s", action)
	}
	topic := ""
	if i := strings.Index(name, "/"); i >= 0 {
		topic = strings.Trim(name[i+1:], "/")
		name = name[:i]
	}
	sk, ok := r.find(name)
	if !ok {
		if avail := r.namesText(); avail != "" {
			return RegisteredSkill{}, nil, fmt.Errorf("unknown skill: %s (available: %s)", name, avail)
		}
		return RegisteredSkill{}, nil, fmt.Errorf("unknown skill: %s", name)
	}
	if topic != "" && action == ActionShow {
		rest = append([]string{topic}, rest...)
	}
	return sk, rest, nil
}

func (r *Registry) handleShow(header bool, rest []string) error {
	sk, rest, err := r.resolve(rest, ActionShow)
	if err != nil {
		return err
	}
	content, err := r.single(sk).loadContent(rest)
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

func (r *Registry) handleVersion(rest []string) error {
	sk, rest, err := r.resolve(rest, ActionVersion)
	if err != nil {
		return err
	}
	return r.single(sk).handleVersion(rest)
}

func (r *Registry) handleInstall(rest []string) error {
	sk, rest, err := r.resolve(rest, ActionInstall)
	if err != nil {
		return err
	}
	return r.single(sk).handleInstall(rest)
}

func (r *Registry) handleUpdateMany(args []string) error {
	usage := strings.TrimSpace(r.UpdateUsage)
	if usage == "" {
		usage = "skills update"
	}
	skills := make([]UpdateSkill, 0, len(r.Skills))
	for _, sk := range r.Skills {
		extras := sk.ExtraFiles
		if extras == nil && sk.TreeFS != nil {
			collected, err := collectTreeSkillFiles(sk.TreeFS)
			if err != nil {
				return err
			}
			extras = collected
		}
		skills = append(skills, UpdateSkill{
			InstallOptions: InstallOptions{
				SkillDirName: sk.Name,
				SkillContent: sk.Content,
				ExtraFiles:   extras,
				Usage:        usage,
			},
			Name: sk.Name,
		})
	}
	return HandleUpdateMany(skills, args)
}
