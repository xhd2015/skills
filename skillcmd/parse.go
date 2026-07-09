package skillcmd

import (
	"fmt"
)

// Action is the skill CLI action selected by flags.
type Action string

const (
	ActionShow    Action = "show"
	ActionInstall Action = "install"
	ActionList    Action = "list"
	ActionHelp    Action = "help"
)

// ParsedArgs is the result of scanning skill argv for action flags.
type ParsedArgs struct {
	Action Action
	Header bool
	Rest   []string
}

// known skill action/mode tokens (lookup avoids switch-on-flag-literal patterns).
var skillArgKind = map[string]string{
	"--show":    "show",
	"--install": "install",
	"--list":    "list",
	"-l":        "list",
	"--header":  "header",
	"-h":        "help",
	"--help":    "help",
}

// ParseSkillArgs scans args for skill action flags only (--show, --install,
// --list / -l, optional --header, -h/--help). Remaining tokens stay in Rest
// (topic path, skill name, install flags like --global).
//
// Help rules (so each skill level is explorable):
//   - -h/--help with no --show/--install/--list → ActionHelp
//   - -h/--help with --show or --list → ActionHelp (skill-level help)
//   - -h/--help with --install only → ActionInstall and --help left in Rest
//     so HandleInstall can print install usage
//
// Exactly one of show/install/list is required when help is not selected for
// the skill-level path; combining --show with --install is an error.
func ParseSkillArgs(args []string) (ParsedArgs, error) {
	var (
		show    bool
		install bool
		list    bool
		header  bool
		help    bool
		rest    []string
	)
	for _, a := range args {
		kind, ok := skillArgKind[a]
		if !ok {
			rest = append(rest, a)
			continue
		}
		switch kind {
		case "show":
			show = true
		case "install":
			install = true
		case "list":
			list = true
		case "header":
			header = true
		case "help":
			help = true
		}
	}

	// Install owns its own --help (targets, --global, etc.).
	if help && install && !show && !list {
		rest = append(rest, "--help")
		return ParsedArgs{
			Action: ActionInstall,
			Rest:   rest,
		}, nil
	}

	// Skill-level help: bare --help, or --help with show/list.
	if help {
		return ParsedArgs{
			Action: ActionHelp,
			Rest:   rest,
		}, nil
	}

	n := 0
	var action Action
	if show {
		n++
		action = ActionShow
	}
	if install {
		n++
		action = ActionInstall
	}
	if list {
		n++
		action = ActionList
	}
	if n == 0 {
		return ParsedArgs{}, fmt.Errorf("expected one of --show, --install, or --list (try --help)")
	}
	if n > 1 {
		if show && install {
			return ParsedArgs{}, fmt.Errorf("cannot combine --show and --install")
		}
		return ParsedArgs{}, fmt.Errorf("expected exactly one of --show, --install, or --list (try --help)")
	}
	if header && action != ActionShow {
		return ParsedArgs{}, fmt.Errorf("--header is only valid with --show")
	}
	return ParsedArgs{
		Action: action,
		Header: header,
		Rest:   rest,
	}, nil
}
