// Package install provides skill installation helpers.
//
// Deprecated: use package github.com/xhd2015/skills/skillcmd instead.
// This package is a thin re-export shim and will be removed in a future version.
package install

import "github.com/xhd2015/skills/skillcmd"

// InstallOptions configures skill installation.
//
// Deprecated: use skillcmd.InstallOptions.
type InstallOptions = skillcmd.InstallOptions

// InstallFile is an extra file written alongside SKILL.md.
//
// Deprecated: use skillcmd.InstallFile.
type InstallFile = skillcmd.InstallFile

// HandleInstall installs a skill to one or more target directories.
//
// Deprecated: use skillcmd.HandleInstall.
func HandleInstall(opts InstallOptions, args []string) error {
	return skillcmd.HandleInstall(opts, args)
}
