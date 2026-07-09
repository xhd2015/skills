package install

import "github.com/xhd2015/skills/skillcmd"

// UpdateSkill describes one entry in a multi-skill registry update batch.
//
// Deprecated: use skillcmd.UpdateSkill.
type UpdateSkill = skillcmd.UpdateSkill

// TargetFlags selects install/update destination directories.
//
// Deprecated: use skillcmd.TargetFlags.
type TargetFlags = skillcmd.TargetFlags

// ResolveTargetDirs maps target flags and optional positional <dir> to paths.
//
// Deprecated: use skillcmd.ResolveTargetDirs.
func ResolveTargetDirs(skillDirName string, tf TargetFlags, args []string) ([]string, error) {
	return skillcmd.ResolveTargetDirs(skillDirName, tf, args)
}

// IsInstalled reports whether SKILL.md exists under dir.
//
// Deprecated: use skillcmd.IsInstalled.
func IsInstalled(dir string) bool {
	return skillcmd.IsInstalled(dir)
}

// InstallTo installs or updates skill content at dir.
//
// Deprecated: use skillcmd.InstallTo.
func InstallTo(dir string, skillContent string, extraFiles []InstallFile, dryRun bool, noOverride bool) error {
	return skillcmd.InstallTo(dir, skillContent, extraFiles, dryRun, noOverride)
}

// HandleUpdate updates one skill at resolved targets that already have SKILL.md.
//
// Deprecated: use skillcmd.HandleUpdate.
func HandleUpdate(opts InstallOptions, args []string) error {
	return skillcmd.HandleUpdate(opts, args)
}

// HandleUpdateMany updates many skills; skips a skill when none of its targets are installed.
//
// Deprecated: use skillcmd.HandleUpdateMany.
func HandleUpdateMany(skills []UpdateSkill, args []string) error {
	return skillcmd.HandleUpdateMany(skills, args)
}
