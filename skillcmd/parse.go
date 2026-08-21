package skillcmd

import (
	"fmt"

	lessflags "github.com/xhd2015/less-flags"
)

// Action is the skill CLI action selected by flags.
type Action string

const (
	ActionShow    Action = "show"
	ActionInstall Action = "install"
	ActionList    Action = "list"
	ActionVersion Action = "version"
	ActionHelp    Action = "help"
)

// ParsedArgs is the result of scanning skill argv for action flags.
type ParsedArgs struct {
	Action Action
	Header bool
	Rest   []string
}

type skillFlagValues struct {
	show    bool
	install bool
	list    bool
	version bool
	header  bool
	help    bool
}

var sharedSkillFlagNames = map[string]struct{}{
	"--show":    {},
	"--install": {},
	"--list":    {},
	"-l":        {},
	"--version": {},
	"--header":  {},
	"-h":        {},
	"--help":    {},
}

// splitSharedSkillFlags keeps install-specific and future flags in Rest rather
// than asking less-flags to reject flags owned by downstream installers.
func splitSharedSkillFlags(args []string) (flagArgs, rest []string) {
	for _, arg := range args {
		if _, ok := sharedSkillFlagNames[arg]; ok {
			flagArgs = append(flagArgs, arg)
			continue
		}
		rest = append(rest, arg)
	}
	return flagArgs, rest
}

func parseSkillFlagValues(args []string) (skillFlagValues, []string, error) {
	flagArgs, rest := splitSharedSkillFlags(args)
	var flags skillFlagValues
	_, err := lessflags.Bool("--show", &flags.show).
		Bool("--install", &flags.install).
		Bool("-l,--list", &flags.list).
		Bool("--version", &flags.version).
		Bool("--header", &flags.header).
		Bool("-h,--help", &flags.help).
		Parse(flagArgs)
	if err != nil {
		return skillFlagValues{}, nil, err
	}
	return flags, rest, nil
}

func selectedSkillAction(flags skillFlagValues) (Action, int) {
	var action Action
	n := 0
	if flags.show {
		n++
		action = ActionShow
	}
	if flags.install {
		n++
		action = ActionInstall
	}
	if flags.list {
		n++
		action = ActionList
	}
	if flags.version {
		n++
		action = ActionVersion
	}
	return action, n
}

// ParseSkillArgs scans args for skill action flags only (--show, --install,
// --list / -l, --version, optional --header, -h/--help). Remaining tokens stay
// in Rest (topic path, skill name, install flags like --global).
//
// Help rules (so each skill level is explorable):
//   - -h/--help with no action → ActionHelp
//   - -h/--help with --show, --list, or --version → ActionHelp
//   - -h/--help with --install only → ActionInstall and --help left in Rest
//     so HandleInstall can print install usage
//
// Exactly one action is required when help is not selected. --header is valid
// only with --show.
func ParseSkillArgs(args []string) (ParsedArgs, error) {
	flags, rest, err := parseSkillFlagValues(args)
	if err != nil {
		return ParsedArgs{}, err
	}
	action, n := selectedSkillAction(flags)
	if n > 1 {
		if flags.show && flags.install && !flags.list && !flags.version {
			return ParsedArgs{}, fmt.Errorf("cannot combine --show and --install")
		}
		return ParsedArgs{}, fmt.Errorf("expected exactly one of --show, --install, --list, or --version (try --help)")
	}
	if flags.header && n == 1 && action != ActionShow {
		return ParsedArgs{}, fmt.Errorf("--header is only valid with --show")
	}

	if flags.help && action == ActionInstall {
		rest = append(rest, "--help")
		return ParsedArgs{Action: ActionInstall, Rest: rest}, nil
	}
	if flags.help {
		return ParsedArgs{Action: ActionHelp, Rest: rest}, nil
	}
	if n == 0 {
		return ParsedArgs{}, fmt.Errorf("expected one of --show, --install, --list, or --version (try --help)")
	}
	if flags.header && action != ActionShow {
		return ParsedArgs{}, fmt.Errorf("--header is only valid with --show")
	}
	return ParsedArgs{Action: action, Header: flags.header, Rest: rest}, nil
}
