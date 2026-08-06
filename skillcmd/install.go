package skillcmd

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/xhd2015/less-gen/flags"
)

type InstallOptions struct {
	SkillDirName string
	// CursorDirName is kept for compatibility with older callers.
	CursorDirName string
	SkillContent  string
	ExtraFiles    []InstallFile
	Usage         string
}

type InstallFile struct {
	Path    string
	Content []byte
}

func HandleInstall(opts InstallOptions, args []string) error {
	usage := strings.TrimSpace(opts.Usage)
	if usage == "" {
		usage = "install"
	}
	skillDirName := skillDirNameFrom(opts)
	var dryRun bool
	var tf TargetFlags
	var noOverride bool
	var force bool
	args, err := flags.Bool("--dry-run", &dryRun).
		Bool("--cursor", &tf.Cursor).
		Bool("--codex", &tf.Codex).
		Bool("--opencode", &tf.Opencode).
		Bool("--general-agents", &tf.GeneralAgents).
		Bool("--no-override", &noOverride).
		Bool("--force", &force).
		Bool("--global", &tf.Global).
		Help("-h,--help", fmt.Sprintf(`
Usage: %s [OPTIONS] [<dir>]

Install skill SKILL.md to a directory.
When no <dir> or target flag is provided, install to .agents/skills/%s.

Options:
  --cursor     Install to .cursor/skills/%s (no dir argument needed)
  --codex      Install to .codex/skills/%s (no dir argument needed)
  --opencode   Install to .opencode/skills/%s (no dir argument needed)
  --general-agents
               Install to .agents/skills/%s (no dir argument needed)
  --global     Install to ~/.<dir>/... instead of current directory's .<dir>/...
  --no-override
               Do not automatically overwrite an existing non-empty directory
  --dry-run    Show what would be created without actually creating anything

Multiple --cursor/--codex/--opencode/--general-agents flags can be combined to install to several locations at once.
`, usage, skillDirName, skillDirName, skillDirName, skillDirName, skillDirName)).Parse(args)
	if err != nil {
		return err
	}
	if force {
		noOverride = false
	}

	dirs, err := ResolveTargetDirs(skillDirName, tf, args)
	if err != nil {
		return err
	}

	for _, dir := range dirs {
		if err := installTo(dir, opts.SkillContent, opts.ExtraFiles, dryRun, noOverride); err != nil {
			return err
		}
	}
	return nil
}

// inventoryAction is one create/update/delete in the install plan.
type inventoryAction struct {
	op      string // "create", "update", or "delete"
	absPath string
	relPath string // relative to skill dir, for sorting
	content []byte // only for create/update
}

// inventoryPlan is the computed create/update/delete set for one skill directory.
type inventoryPlan struct {
	absDir   string
	exists   bool
	nonEmpty bool
	actions  []inventoryAction
}

// planInventory compares desired skill content against on-disk files without writing.
func planInventory(dir string, skillContent string, extraFiles []InstallFile) (*inventoryPlan, error) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, fmt.Errorf("resolve path: %w", err)
	}
	planned, err := buildPlannedFiles(absDir, skillContent, extraFiles)
	if err != nil {
		return nil, err
	}

	entries, readErr := os.ReadDir(absDir)
	if readErr != nil && !os.IsNotExist(readErr) {
		return nil, fmt.Errorf("read directory %s: %w", absDir, readErr)
	}
	exists := readErr == nil
	nonEmpty := exists && len(entries) > 0

	onDisk, err := listRegularFiles(absDir, exists)
	if err != nil {
		return nil, err
	}
	actions, err := computeInventoryActions(absDir, planned, onDisk)
	if err != nil {
		return nil, err
	}
	return &inventoryPlan{
		absDir:   absDir,
		exists:   exists,
		nonEmpty: nonEmpty,
		actions:  actions,
	}, nil
}

func installTo(dir string, skillContent string, extraFiles []InstallFile, dryRun bool, noOverride bool) error {
	plan, err := planInventory(dir, skillContent, extraFiles)
	if err != nil {
		return err
	}
	if len(plan.actions) == 0 {
		printUpToDate(plan.absDir, dryRun)
		return nil
	}
	if !dryRun && plan.nonEmpty && noOverride && !confirmOverwrite(plan.absDir) {
		fmt.Println("Aborted.")
		return nil
	}
	header := installHeader(plan.absDir, plan.exists)
	if dryRun {
		printInventoryPlan(header, plan.actions, true)
		return nil
	}
	if err := applyInventoryActions(plan.absDir, plan.actions); err != nil {
		return err
	}
	printInventoryPlan(header, plan.actions, false)
	return nil
}

func printUpToDate(absDir string, dryRun bool) {
	if dryRun {
		fmt.Printf("[dry-run] Skill is up to date: %s\n", absDir)
		return
	}
	fmt.Printf("Skill is up to date: %s\n", absDir)
}

func installHeader(absDir string, exists bool) string {
	if !exists {
		return fmt.Sprintf("Installed skill to: %s", absDir)
	}
	return fmt.Sprintf("Update skill at %s", absDir)
}

func printInventoryPlan(header string, actions []inventoryAction, dryRun bool) {
	if dryRun {
		fmt.Printf("[dry-run] %s\n", header)
		for _, a := range actions {
			fmt.Printf("[dry-run]   %s: %s\n", a.op, a.absPath)
		}
		return
	}
	fmt.Printf("%s\n", header)
	for _, a := range actions {
		fmt.Printf("  %s: %s\n", a.op, a.absPath)
	}
}

// plannedFile is one path in the desired install set.
type plannedFile struct {
	relPath string // slash or OS relative path under skill dir
	absPath string
	content []byte
}

func buildPlannedFiles(absDir string, skillContent string, extraFiles []InstallFile) ([]plannedFile, error) {
	extras, err := resolveInstallFiles(absDir, extraFiles)
	if err != nil {
		return nil, err
	}
	planned := make([]plannedFile, 0, 1+len(extras))
	planned = append(planned, plannedFile{
		relPath: "SKILL.md",
		absPath: filepath.Join(absDir, "SKILL.md"),
		content: []byte(skillContent),
	})
	for _, f := range extras {
		rel, err := filepath.Rel(absDir, f.dest)
		if err != nil {
			rel = f.dest
		}
		planned = append(planned, plannedFile{
			relPath: rel,
			absPath: f.dest,
			content: f.content,
		})
	}
	return planned, nil
}

// listRegularFiles returns a map of relPath → absPath for all regular files under absDir.
// If the directory does not exist, returns an empty map.
func listRegularFiles(absDir string, exists bool) (map[string]string, error) {
	out := make(map[string]string)
	if !exists {
		return out, nil
	}
	err := filepath.WalkDir(absDir, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		// Only inventory regular files (skip symlinks to dirs etc. via type check).
		info, infoErr := d.Info()
		if infoErr != nil {
			return infoErr
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		rel, relErr := filepath.Rel(absDir, p)
		if relErr != nil {
			return relErr
		}
		if rel == "." {
			return nil
		}
		out[rel] = p
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("scan skill directory %s: %w", absDir, err)
	}
	return out, nil
}

func computeInventoryActions(absDir string, planned []plannedFile, onDisk map[string]string) ([]inventoryAction, error) {
	plannedRel := make(map[string]plannedFile, len(planned))
	for _, p := range planned {
		// Normalize rel path for map key consistency with WalkDir
		rel := filepath.Clean(p.relPath)
		plannedRel[rel] = p
	}

	var creates, updates, deletes []inventoryAction

	for rel, p := range plannedRel {
		absPath := p.absPath
		if absPath == "" {
			absPath = filepath.Join(absDir, rel)
		}
		_, diskSame, err := sameFileMD5(absPath, p.content)
		if err != nil {
			return nil, err
		}
		if _, existsOnDisk := onDisk[rel]; !existsOnDisk {
			// File missing (or not a regular file we scanned)
			if !diskSame {
				creates = append(creates, inventoryAction{
					op:      "create",
					absPath: absPath,
					relPath: rel,
					content: p.content,
				})
			}
			continue
		}
		if !diskSame {
			updates = append(updates, inventoryAction{
				op:      "update",
				absPath: absPath,
				relPath: rel,
				content: p.content,
			})
		}
	}

	for rel, absPath := range onDisk {
		if _, inPlan := plannedRel[rel]; inPlan {
			continue
		}
		deletes = append(deletes, inventoryAction{
			op:      "delete",
			absPath: absPath,
			relPath: rel,
		})
	}

	sort.Slice(creates, func(i, j int) bool { return creates[i].relPath < creates[j].relPath })
	sort.Slice(updates, func(i, j int) bool { return updates[i].relPath < updates[j].relPath })
	sort.Slice(deletes, func(i, j int) bool { return deletes[i].relPath < deletes[j].relPath })

	actions := make([]inventoryAction, 0, len(creates)+len(updates)+len(deletes))
	actions = append(actions, creates...)
	actions = append(actions, updates...)
	actions = append(actions, deletes...)
	return actions, nil
}

func applyInventoryActions(absDir string, actions []inventoryAction) error {
	// Ensure skill root exists when creating files.
	if err := os.MkdirAll(absDir, 0755); err != nil {
		return fmt.Errorf("create directory %s: %w", absDir, err)
	}

	deletedParents := make(map[string]struct{})
	for _, a := range actions {
		switch a.op {
		case "create", "update":
			if err := os.MkdirAll(filepath.Dir(a.absPath), 0755); err != nil {
				return fmt.Errorf("create directory %s: %w", filepath.Dir(a.absPath), err)
			}
			if err := os.WriteFile(a.absPath, a.content, 0644); err != nil {
				return fmt.Errorf("write %s: %w", a.absPath, err)
			}
		case "delete":
			if err := os.Remove(a.absPath); err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("delete %s: %w", a.absPath, err)
			}
			// Track parent for empty-dir cleanup
			parent := filepath.Dir(a.absPath)
			if parent != absDir && strings.HasPrefix(parent, absDir+string(os.PathSeparator)) {
				deletedParents[parent] = struct{}{}
			}
		}
	}

	// Best-effort: remove empty directories under skill root (deepest first).
	if len(deletedParents) > 0 {
		parents := make([]string, 0, len(deletedParents))
		for p := range deletedParents {
			parents = append(parents, p)
		}
		// Also walk up ancestors under absDir
		all := make(map[string]struct{})
		for _, p := range parents {
			cur := p
			for cur != absDir && strings.HasPrefix(cur, absDir+string(os.PathSeparator)) {
				all[cur] = struct{}{}
				next := filepath.Dir(cur)
				if next == cur {
					break
				}
				cur = next
			}
		}
		dirs := make([]string, 0, len(all))
		for d := range all {
			dirs = append(dirs, d)
		}
		// deepest first
		sort.Slice(dirs, func(i, j int) bool {
			return strings.Count(dirs[i], string(os.PathSeparator)) > strings.Count(dirs[j], string(os.PathSeparator))
		})
		for _, d := range dirs {
			// Remove only if empty; ignore errors (best-effort).
			_ = os.Remove(d)
		}
	}
	return nil
}

type resolvedInstallFile struct {
	dest    string
	content []byte
}

func resolveInstallFiles(absDir string, extraFiles []InstallFile) ([]resolvedInstallFile, error) {
	files := make([]resolvedInstallFile, 0, len(extraFiles))
	for _, f := range extraFiles {
		rel := filepath.Clean(filepath.FromSlash(f.Path))
		if rel == "." || filepath.IsAbs(rel) || rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
			return nil, fmt.Errorf("invalid install file path: %s", f.Path)
		}
		if rel == "SKILL.md" {
			return nil, fmt.Errorf("extra install file cannot replace SKILL.md")
		}
		files = append(files, resolvedInstallFile{
			dest:    filepath.Join(absDir, rel),
			content: f.Content,
		})
	}
	return files, nil
}

func sameFileMD5(path string, content []byte) (string, bool, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", false, nil
		}
		return "", false, fmt.Errorf("open existing file %s: %w", path, err)
	}
	defer f.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, f); err != nil {
		return "", false, fmt.Errorf("read existing file %s: %w", path, err)
	}
	currentMD5 := fmt.Sprintf("%x", hash.Sum(nil))
	return currentMD5, currentMD5 == md5Hex(content), nil
}

func md5Hex(content []byte) string {
	return fmt.Sprintf("%x", md5.Sum(content))
}

func joinSkillDir(toolDir, skillDirName string) string {
	return filepath.Join(toolDir, "skills", skillDirName)
}

func userHomeDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get home dir: %w", err)
	}
	return homeDir, nil
}

func joinHomeIfRelative(homeDir, dir string) string {
	if filepath.IsAbs(dir) {
		return dir
	}
	return filepath.Join(homeDir, dir)
}

func absPath(dir string) (string, error) {
	return filepath.Abs(dir)
}

func skillMDPathStat(absDir string) (os.FileInfo, error) {
	return os.Stat(filepath.Join(absDir, "SKILL.md"))
}

func confirmOverwrite(dir string) bool {
	f, _ := os.Stdin.Stat()
	if f == nil || (f.Mode()&os.ModeCharDevice) == 0 {
		return false
	}
	fmt.Printf("Directory %s is not empty. Overwrite? [y/N] ", dir)
	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(strings.ToLower(answer))
	return answer == "y" || answer == "yes"
}
