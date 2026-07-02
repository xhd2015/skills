package install

import (
	"errors"
	"fmt"
	"strings"

	"github.com/xhd2015/less-gen/flags"
)

// UpdateSkill describes one entry in a multi-skill registry.
type UpdateSkill struct {
	InstallOptions
	// Name is the CLI alias shown in messages; falls back to SkillDirName when empty.
	Name string
}

// TargetFlags selects install/update destination directories (without --dry-run).
type TargetFlags struct {
	Cursor        bool
	Codex         bool
	Opencode      bool
	GeneralAgents bool
	Global        bool
}

// ResolveTargetDirs maps target flags and optional positional <dir> to filesystem paths.
func ResolveTargetDirs(skillDirName string, tf TargetFlags, args []string) ([]string, error) {
	if skillDirName == "" {
		return nil, fmt.Errorf("skill directory name is required")
	}

	var dirs []string
	if tf.Cursor {
		dirs = append(dirs, joinSkillDir(".cursor", skillDirName))
	}
	if tf.Codex {
		dirs = append(dirs, joinSkillDir(".codex", skillDirName))
	}
	if tf.Opencode {
		dirs = append(dirs, joinSkillDir(".opencode", skillDirName))
	}
	if tf.GeneralAgents {
		dirs = append(dirs, joinSkillDir(".agents", skillDirName))
	}
	if len(dirs) == 0 {
		if len(args) > 0 {
			dirs = append(dirs, args[0])
		} else {
			dirs = append(dirs, joinSkillDir(".agents", skillDirName))
		}
	}

	if tf.Global {
		homeDir, err := userHomeDir()
		if err != nil {
			return nil, err
		}
		for i, dir := range dirs {
			dirs[i] = joinHomeIfRelative(homeDir, dir)
		}
	}
	return dirs, nil
}

// IsInstalled reports whether SKILL.md exists under dir.
func IsInstalled(dir string) bool {
	absDir, err := absPath(dir)
	if err != nil {
		return false
	}
	_, err = skillMDPathStat(absDir)
	return err == nil
}

// InstallTo installs or updates skill content at dir (exported installTo).
func InstallTo(dir string, skillContent string, extraFiles []InstallFile, dryRun bool, noOverride bool) error {
	return installTo(dir, skillContent, extraFiles, dryRun, noOverride)
}

// HandleUpdate updates one skill at resolved targets that already have SKILL.md.
func HandleUpdate(opts InstallOptions, args []string) error {
	tf, dryRun, args, err := parseUpdateFlags(opts, args)
	if err != nil {
		if errors.Is(err, flags.ErrHelp) {
			return nil
		}
		return err
	}
	skillDirName := skillDirNameFrom(opts)
	dirs, err := ResolveTargetDirs(skillDirName, tf, args)
	if err != nil {
		return err
	}
	for _, dir := range dirs {
		if !IsInstalled(dir) {
			continue
		}
		if err := InstallTo(dir, opts.SkillContent, opts.ExtraFiles, dryRun, false); err != nil {
			return err
		}
	}
	return nil
}

// HandleUpdateMany updates many skills; skips a skill when none of its targets are installed.
func HandleUpdateMany(skills []UpdateSkill, args []string) error {
	if len(skills) == 0 {
		return nil
	}
	usage := strings.TrimSpace(skills[0].Usage)
	if usage == "" {
		usage = "skills update"
	}
	tf, dryRun, args, err := parseUpdateFlags(InstallOptions{Usage: usage}, args)
	if err != nil {
		if errors.Is(err, flags.ErrHelp) {
			return nil
		}
		return err
	}
	for _, skill := range skills {
		skillDirName := skillDirNameFrom(skill.InstallOptions)
		dirs, err := ResolveTargetDirs(skillDirName, tf, args)
		if err != nil {
			return err
		}
		var installed []string
		for _, dir := range dirs {
			if IsInstalled(dir) {
				installed = append(installed, dir)
			}
		}
		if len(installed) == 0 {
			fmt.Println("skill not installed: " + updateSkillDisplayName(skill))
			continue
		}
		for _, dir := range installed {
			if err := InstallTo(dir, skill.SkillContent, skill.ExtraFiles, dryRun, false); err != nil {
				return err
			}
		}
		fmt.Println(skillDirName)
	}
	return nil
}

func updateSkillDisplayName(skill UpdateSkill) string {
	if skill.Name != "" {
		return skill.Name
	}
	return skillDirNameFrom(skill.InstallOptions)
}

func parseUpdateFlags(opts InstallOptions, args []string) (TargetFlags, bool, []string, error) {
	usage := strings.TrimSpace(opts.Usage)
	if usage == "" {
		usage = "update"
	}
	skillDirName := skillDirNameFrom(opts)
	var tf TargetFlags
	var dryRun bool
	args, err := flags.Bool("--dry-run", &dryRun).
		Bool("--cursor", &tf.Cursor).
		Bool("--codex", &tf.Codex).
		Bool("--opencode", &tf.Opencode).
		Bool("--general-agents", &tf.GeneralAgents).
		Bool("--global", &tf.Global).
		Help("-h,--help", fmt.Sprintf(`
Usage: %s [OPTIONS] [<dir>]

Update an already-installed skill (SKILL.md must exist at target paths).
When no <dir> or target flag is provided, update .agents/skills/%s.

Options:
  --cursor     Update .cursor/skills/%s when installed
  --codex      Update .codex/skills/%s when installed
  --opencode   Update .opencode/skills/%s when installed
  --general-agents
               Update .agents/skills/%s when installed
  --global     Use ~/.<dir>/... instead of current directory's .<dir>/...
  --dry-run    Show what would change without writing files

Multiple --cursor/--codex/--opencode/--general-agents flags can be combined.
`, usage, skillDirName, skillDirName, skillDirName, skillDirName, skillDirName)).HelpNoExit().Parse(args)
	if err != nil {
		return TargetFlags{}, false, nil, err
	}
	return tf, dryRun, args, nil
}

func skillDirNameFrom(opts InstallOptions) string {
	if opts.SkillDirName != "" {
		return opts.SkillDirName
	}
	return opts.CursorDirName
}
