package skillcmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/xhd2015/dot-pkgs/go-pkgs/terminal/color"
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

// updateSkillResult is the polished per-skill report for update handlers.
type updateSkillResult struct {
	name    string
	status  string // "up to date" | "updated" | "would update" | "not installed"
	actions []inventoryAction
}

// HandleUpdate updates one skill at resolved targets that already have SKILL.md.
// Silent when nothing is installed. Prints polished status + optional file lines.
func HandleUpdate(opts InstallOptions, args []string) error {
	return handleUpdate(os.Stdout, opts, args)
}

func handleUpdate(w io.Writer, opts InstallOptions, args []string) error {
	tf, dryRun, colorMode, args, err := parseUpdateFlags(opts, args)
	if err != nil {
		if errors.Is(err, flags.ErrHelp) {
			return nil
		}
		return err
	}
	style := color.Style{Enabled: color.EnabledFor(colorMode, w)}
	skillDirName := skillDirNameFrom(opts)
	dirs, err := ResolveTargetDirs(skillDirName, tf, args)
	if err != nil {
		return err
	}
	name := skillDirName
	if opts.SkillDirName != "" {
		name = opts.SkillDirName
	}
	// Prefer CursorDirName only as skill dir; display uses skillDirName.
	result, err := runUpdateSkill(name, opts, dirs, dryRun)
	if err != nil {
		return err
	}
	if result.status == "not installed" {
		// Single update: silent when nothing installed.
		return nil
	}
	printUpdateSkillResult(w, style, result)
	return nil
}

// HandleUpdateMany updates many skills; reports not installed for skills with no targets.
func HandleUpdateMany(skills []UpdateSkill, args []string) error {
	return handleUpdateMany(os.Stdout, skills, args)
}

func handleUpdateMany(w io.Writer, skills []UpdateSkill, args []string) error {
	if len(skills) == 0 {
		return nil
	}
	usage := strings.TrimSpace(skills[0].Usage)
	if usage == "" {
		usage = "skills update"
	}
	tf, dryRun, colorMode, args, err := parseUpdateFlags(InstallOptions{Usage: usage}, args)
	if err != nil {
		if errors.Is(err, flags.ErrHelp) {
			return nil
		}
		return err
	}
	style := color.Style{Enabled: color.EnabledFor(colorMode, w)}

	var nUpdated, nWouldUpdate, nUpToDate, nNotInstalled int
	for _, skill := range skills {
		skillDirName := skillDirNameFrom(skill.InstallOptions)
		dirs, err := ResolveTargetDirs(skillDirName, tf, args)
		if err != nil {
			return err
		}
		name := updateSkillDisplayName(skill)
		result, err := runUpdateSkill(name, skill.InstallOptions, dirs, dryRun)
		if err != nil {
			return err
		}
		printUpdateSkillResult(w, style, result)
		switch result.status {
		case "updated":
			nUpdated++
		case "would update":
			nWouldUpdate++
		case "up to date":
			nUpToDate++
		case "not installed":
			nNotInstalled++
		}
	}
	fmt.Fprintln(w)
	summary := formatUpdateSummary(nUpdated, nWouldUpdate, nUpToDate, nNotInstalled, dryRun)
	fmt.Fprintln(w, style.Gray(summary))
	return nil
}

// runUpdateSkill plans (and optionally applies) inventory for installed targets only.
func runUpdateSkill(name string, opts InstallOptions, dirs []string, dryRun bool) (updateSkillResult, error) {
	var installed []string
	for _, dir := range dirs {
		if IsInstalled(dir) {
			installed = append(installed, dir)
		}
	}
	if len(installed) == 0 {
		return updateSkillResult{name: name, status: "not installed"}, nil
	}

	var all []inventoryAction
	for _, dir := range installed {
		plan, err := planInventory(dir, opts.SkillContent, opts.ExtraFiles)
		if err != nil {
			return updateSkillResult{}, err
		}
		if len(plan.actions) == 0 {
			continue
		}
		if !dryRun {
			if err := applyInventoryActions(plan.absDir, plan.actions); err != nil {
				return updateSkillResult{}, err
			}
		}
		all = append(all, plan.actions...)
	}

	if len(all) == 0 {
		return updateSkillResult{name: name, status: "up to date"}, nil
	}
	// Stable create → update → delete across multi-target aggregation.
	all = sortInventoryActions(all)
	status := "updated"
	if dryRun {
		status = "would update"
	}
	return updateSkillResult{name: name, status: status, actions: all}, nil
}

func sortInventoryActions(actions []inventoryAction) []inventoryAction {
	opOrder := map[string]int{"create": 0, "update": 1, "delete": 2}
	out := append([]inventoryAction(nil), actions...)
	sort.SliceStable(out, func(i, j int) bool {
		oi, oj := opOrder[out[i].op], opOrder[out[j].op]
		if oi != oj {
			return oi < oj
		}
		if out[i].relPath != out[j].relPath {
			return out[i].relPath < out[j].relPath
		}
		return out[i].absPath < out[j].absPath
	})
	return out
}

func printUpdateSkillResult(w io.Writer, style color.Style, r updateSkillResult) {
	status := colorStatus(style, r.status)
	if len(r.actions) == 0 {
		fmt.Fprintf(w, "%s  %s\n", r.name, status)
		return
	}
	counts := formatOpCounts(r.actions)
	// Counts and parentheses stay plain; only the status token is tinted.
	fmt.Fprintf(w, "%s  %s  (%s)\n", r.name, status, counts)
	for _, a := range r.actions {
		fmt.Fprintf(w, "  %s  %s\n", style.Gray(a.op), a.absPath)
	}
}

// formatOpCounts builds "1 create, 6 update" (only non-zero ops; create, update, delete order).
func formatOpCounts(actions []inventoryAction) string {
	var creates, updates, deletes int
	for _, a := range actions {
		switch a.op {
		case "create":
			creates++
		case "update":
			updates++
		case "delete":
			deletes++
		}
	}
	var parts []string
	if creates > 0 {
		parts = append(parts, fmt.Sprintf("%d create", creates))
	}
	if updates > 0 {
		parts = append(parts, fmt.Sprintf("%d update", updates))
	}
	if deletes > 0 {
		parts = append(parts, fmt.Sprintf("%d delete", deletes))
	}
	return strings.Join(parts, ", ")
}

// formatUpdateSummary builds the trailing batch count line.
func formatUpdateSummary(updated, wouldUpdate, upToDate, notInstalled int, dryRun bool) string {
	if dryRun {
		return fmt.Sprintf("%d updated · %d would update · %d up to date · %d not installed  [dry-run]",
			updated, wouldUpdate, upToDate, notInstalled)
	}
	return fmt.Sprintf("%d updated · %d up to date · %d not installed",
		updated, upToDate, notInstalled)
}

func updateSkillDisplayName(skill UpdateSkill) string {
	if skill.Name != "" {
		return skill.Name
	}
	return skillDirNameFrom(skill.InstallOptions)
}

func parseUpdateFlags(opts InstallOptions, args []string) (TargetFlags, bool, color.Mode, []string, error) {
	usage := strings.TrimSpace(opts.Usage)
	if usage == "" {
		usage = "update"
	}
	skillDirName := skillDirNameFrom(opts)
	var tf TargetFlags
	var dryRun bool
	var colorFlag, noColorFlag bool
	args, err := flags.Bool("--dry-run", &dryRun).
		Bool("--cursor", &tf.Cursor).
		Bool("--codex", &tf.Codex).
		Bool("--opencode", &tf.Opencode).
		Bool("--general-agents", &tf.GeneralAgents).
		Bool("--global", &tf.Global).
		Bool("--color", &colorFlag).
		Bool("--no-color", &noColorFlag).
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
  --color      force ANSI color on (even when stdout is not a TTY)
  --no-color   force ANSI color off

Multiple --cursor/--codex/--opencode/--general-agents flags can be combined.
`, usage, skillDirName, skillDirName, skillDirName, skillDirName, skillDirName)).HelpNoExit().Parse(args)
	if err != nil {
		return TargetFlags{}, false, color.Auto, nil, err
	}
	mode, err := color.ModeFromFlags(colorFlag, noColorFlag)
	if err != nil {
		return TargetFlags{}, false, color.Auto, nil, err
	}
	return tf, dryRun, mode, args, nil
}

func skillDirNameFrom(opts InstallOptions) string {
	if opts.SkillDirName != "" {
		return opts.SkillDirName
	}
	return opts.CursorDirName
}
