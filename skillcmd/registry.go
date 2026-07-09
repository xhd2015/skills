package skillcmd

import (
	"fmt"
)

// RegisteredSkill is one entry in a multi-skill registry.
type RegisteredSkill struct {
	Name        string
	Description string
	Content     string
	ExtraFiles  []InstallFile
}

// Registry is an ordered list of registered skills for multi-skill hosts.
type Registry struct {
	Skills []RegisteredSkill
}

// HandleSkill implements skill flag actions: --list | --show|--install name.
func (r *Registry) HandleSkill(args []string) error {
	parsed, err := ParseSkillArgs(args)
	if err != nil {
		return err
	}
	switch parsed.Action {
	case ActionList:
		return r.listSkills()
	case ActionShow:
		return r.handleShow(parsed.Header, parsed.Rest)
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
	if args[0] == "update" {
		return r.handleUpdateMany(args[1:])
	}
	return r.HandleSkill(args)
}

func (r *Registry) listSkills() error {
	for _, sk := range r.Skills {
		if sk.Description != "" {
			fmt.Printf("%s - %s\n", sk.Name, sk.Description)
		} else {
			fmt.Println(sk.Name)
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

func (r *Registry) handleShow(header bool, rest []string) error {
	if len(rest) == 0 {
		return fmt.Errorf("expected skill name for --show")
	}
	name := rest[0]
	if len(rest) > 1 {
		return fmt.Errorf("unexpected arguments: %v", rest[1:])
	}
	sk, ok := r.find(name)
	if !ok {
		return fmt.Errorf("unknown skill: %s", name)
	}
	if header {
		out, err := FormatHeaderWithDelimiters(sk.Content)
		if err != nil {
			return err
		}
		fmt.Print(out)
		return nil
	}
	fmt.Print(sk.Content)
	return nil
}

func (r *Registry) handleInstall(rest []string) error {
	if len(rest) == 0 {
		return fmt.Errorf("expected skill name for --install")
	}
	name := rest[0]
	sk, ok := r.find(name)
	if !ok {
		return fmt.Errorf("unknown skill: %s", name)
	}
	return HandleInstall(InstallOptions{
		SkillDirName: sk.Name,
		SkillContent: sk.Content,
		ExtraFiles:   sk.ExtraFiles,
		Usage:        "skill --install " + sk.Name,
	}, rest[1:])
}

func (r *Registry) handleUpdateMany(args []string) error {
	skills := make([]UpdateSkill, 0, len(r.Skills))
	for _, sk := range r.Skills {
		usage := "skills update"
		skills = append(skills, UpdateSkill{
			InstallOptions: InstallOptions{
				SkillDirName: sk.Name,
				SkillContent: sk.Content,
				ExtraFiles:   sk.ExtraFiles,
				Usage:        usage,
			},
			Name: sk.Name,
		})
	}
	return HandleUpdateMany(skills, args)
}
