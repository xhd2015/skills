package main

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/xhd2015/skills/cmd/go-best-practice/vet"
	"github.com/xhd2015/skills/skillcmd"
)

//go:embed SKILL.md
var skillTemplate string

// Nested topic directories (path/SKILL.md layout).
//
//go:embed cmd-exec
//go:embed flags-parsing
//go:embed kool-create
//go:embed skill-cli
var skillTreeFS embed.FS

const skillName = "go-best-practice"

const help = `
Usage: go-best-practice <command> [ARGS]
       go-best-practice skill --show [<topic>[/<sub-topic>[/...]]]
       go-best-practice skill <topic-path> --show

Commands:
  install [<dir>]          Install SKILL.md + nested topics (alias of skill --install)
  skill --show [--header] [path]
                           Show root SKILL.md or nested path/SKILL.md
  skill --install [<dir>]  Install SKILL.md + nested topics to a directory
  skill --list             Print the skill name
  topics                   List all available top-level topics
  vet [flags] [dirs]       Check codebase for best-practice violations

Topics are nested directories with SKILL.md. Address a nested topic with a
slash-separated path, e.g. "flags-parsing/types", via skill --show.

Run go-best-practice skill --help for skill subcommand options.
Run go-best-practice skill --install --help for install flags.

Options:
  -h, --help    Show this help message
`

const skillHelp = `Usage: go-best-practice skill --show [--header] [<topic-path>]
       go-best-practice skill <topic-path> --show [--header]
       go-best-practice skill --install [OPTIONS] [<dir>]
       go-best-practice skill --list

Show the root SKILL.md index or a nested topic (path/SKILL.md).
Install copies SKILL.md and nested topics into agent skill directories.
List prints the skill name (go-best-practice).

Examples:
  go-best-practice skill --show
  go-best-practice skill --show flags-parsing/types
  go-best-practice skill flags-parsing/types --show
  go-best-practice skill --install --dry-run
  go-best-practice skill --install --help

Options:
  --show [--header] [path]   Print skill or topic content (header-only with --header)
  --install [OPTIONS] [dir]  Install skill files (see --install --help)
  --list                     Print skill directory name
  -h, --help                 Show this help message
`

func main() {
	if err := handle(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func handle(args []string) error {
	if len(args) == 0 {
		fmt.Print(help)
		fmt.Println()
		return printTopicIndex()
	}

	switch args[0] {
	case "-h", "--help":
		fmt.Print(help)
		fmt.Println()
		return printTopicIndex()
	case "install":
		// top-level install alias → skill --install
		return singleSkill().Handle(append([]string{"--install"}, args[1:]...))
	case "topics", "list":
		return printTopicIndex()
	case "skill":
		return singleSkill().Handle(args[1:])
	case "vet":
		return vet.Run(args[1:])
	default:
		return fmt.Errorf("unknown command: %s (use skill --show <topic> or run `go-best-practice topics`)", args[0])
	}
}

func singleSkill() *skillcmd.SingleSkill {
	return &skillcmd.SingleSkill{
		Name:        skillName,
		RootContent: skillTemplate,
		TreeFS:      skillTreeFS,
		Usage:       "go-best-practice skill --install",
		Help:        skillHelp,
	}
}

// listTopics returns all topic paths found under skillTreeFS (directories
// that contain SKILL.md), slash-separated, e.g. "flags-parsing",
// "flags-parsing/types".
func listTopics() ([]string, error) {
	var topics []string
	err := fs.WalkDir(skillTreeFS, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if path.Base(p) != "SKILL.md" {
			return nil
		}
		p = path.Clean(p)
		if p == "SKILL.md" || p == "./SKILL.md" {
			return nil
		}
		rel := strings.TrimSuffix(p, "/SKILL.md")
		rel = strings.TrimSuffix(rel, "SKILL.md")
		rel = strings.Trim(rel, "/")
		if rel == "" {
			return nil
		}
		topics = append(topics, rel)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk embedded skill tree: %w", err)
	}
	sort.Strings(topics)
	return topics, nil
}

func printTopicIndex() error {
	topics, err := listTopics()
	if err != nil {
		return err
	}
	fmt.Println("Available topics:")
	for _, t := range topics {
		depth := strings.Count(t, "/")
		indent := strings.Repeat("  ", depth)
		label := t
		if idx := strings.LastIndex(t, "/"); idx >= 0 {
			label = t[idx+1:]
		}
		fmt.Printf("  %s- %s\n", indent, label)
	}
	return nil
}

// readTopic resolves a slash-separated topic path against nested SKILL.md files.
func readTopic(topicPath string) (string, bool, error) {
	topicPath = strings.Trim(topicPath, "/")
	if topicPath == "" {
		return "", false, nil
	}
	segments := strings.Split(topicPath, "/")
	if err := validateSegments(segments); err != nil {
		return "", false, err
	}
	embedPath := path.Join(topicPath, "SKILL.md")
	data, err := skillTreeFS.ReadFile(embedPath)
	if err == nil {
		return string(data), true, nil
	}
	if os.IsNotExist(err) {
		return "", false, nil
	}
	return "", false, fmt.Errorf("read topic %s: %w", topicPath, err)
}

func validateSegments(segments []string) error {
	for _, s := range segments {
		if s == "" || s == "." || s == ".." {
			return fmt.Errorf("invalid topic path segment: %q", s)
		}
	}
	return nil
}

// collectTopicFiles returns nested SKILL.md files for install ExtraFiles.
func collectTopicFiles() ([]skillcmd.InstallFile, error) {
	var files []skillcmd.InstallFile
	err := fs.WalkDir(skillTreeFS, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		p = path.Clean(p)
		if path.Base(p) != "SKILL.md" || p == "SKILL.md" || p == "./SKILL.md" {
			return nil
		}
		data, err := skillTreeFS.ReadFile(p)
		if err != nil {
			return err
		}
		files = append(files, skillcmd.InstallFile{Path: p, Content: data})
		return nil
	})
	return files, err
}
