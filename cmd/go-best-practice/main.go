package main

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/xhd2015/skills/cmd/go-best-practice/vet"
	"github.com/xhd2015/skills/install"
	"github.com/xhd2015/skills/skill_file"
)

//go:embed SKILL.md
var skillTemplate string

//go:embed topics
var topicsFS embed.FS

const topicsDir = "topics"

const help = `
Usage: go-best-practice <command> [ARGS]
       go-best-practice <topic>[/<sub-topic>[/...]]

Commands:
  install [<dir>]          Install SKILL.md + topics to a directory (or use --cursor)
  skill show               Show the content of SKILL.md
  skill install [<dir>]    Install SKILL.md + topics to a directory
  topics                   List all available top-level topics
  vet [flags] [dirs] Check codebase for best-practice violations
  <topic-path>       Print the detailed content for a topic or sub-topic

Topics are organized hierarchically. Address a nested topic with a
slash-separated path, e.g. "flags-parsing/types".

Options:
  -h, --help    Show this help message
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
		return handleInstall(args[1:])
	case "topics", "list":
		return printTopicIndex()
	case "skill":
		return handleSkill(args[1:])
	case "vet":
		return vet.Run(args[1:])
	}

	content, ok, err := readTopic(args[0])
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("unknown command or topic: %s (run `go-best-practice topics` to list available topics)", args[0])
	}
	fmt.Print(content)
	if !strings.HasSuffix(content, "\n") {
		fmt.Println()
	}
	return nil
}

// listTopics returns all topic paths found under topicsDir, relative
// to topicsDir and slash-separated. E.g. "flags-parsing",
// "flags-parsing/types". Each returned path corresponds to a
// reachable <path>.md file in the embedded FS.
func listTopics() ([]string, error) {
	var topics []string
	err := fs.WalkDir(topicsFS, topicsDir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(p, ".md") {
			return nil
		}
		rel, err := filepath.Rel(topicsDir, p)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		rel = strings.TrimSuffix(rel, ".md")
		topics = append(topics, rel)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk embedded topics: %w", err)
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

// readTopic resolves a slash-separated topic path like
// "flags-quickstart/types" against the embedded topics tree. A path
// segment is matched either as `<segment>.md` (leaf) or as
// `<segment>/` (directory) where the next segment is searched.
func readTopic(topicPath string) (string, bool, error) {
	topicPath = strings.Trim(topicPath, "/")
	if topicPath == "" {
		return "", false, nil
	}

	segments := strings.Split(topicPath, "/")
	if err := validateSegments(segments); err != nil {
		return "", false, err
	}

	embedPath := path.Join(append([]string{topicsDir}, segments...)...) + ".md"
	data, err := topicsFS.ReadFile(embedPath)
	if err == nil {
		return string(data), true, nil
	}
	if !os.IsNotExist(err) {
		return "", false, fmt.Errorf("read topic %s: %w", topicPath, err)
	}
	return "", false, nil
}

func validateSegments(segments []string) error {
	for _, s := range segments {
		if s == "" || s == "." || s == ".." {
			return fmt.Errorf("invalid topic path segment: %q", s)
		}
	}
	return nil
}

func handleSkill(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("unknown skill sub-command: expected `skill show` or `skill install`")
	}
	switch args[0] {
	case "show":
		rest := args[1:]
		headerOnly := false
		if len(rest) > 0 && rest[0] == "--header" {
			headerOnly = true
			rest = rest[1:]
		}
		if len(rest) > 0 {
			return fmt.Errorf("unknown skill show option: %s", rest[0])
		}
		if headerOnly {
			out, err := skill_file.FormatHeaderWithDelimiters(skillTemplate)
			if err != nil {
				return err
			}
			fmt.Print(out)
			return nil
		}
		fmt.Print(skillTemplate)
		return nil
	case "install":
		return handleInstall(args[1:])
	default:
		return fmt.Errorf("unknown skill sub-command: %s (expected `skill show` or `skill install`)", args[0])
	}
}

func handleInstall(args []string) error {
	files, err := collectTopicFiles()
	if err != nil {
		return err
	}
	return install.HandleInstall(install.InstallOptions{
		SkillDirName: "go-best-practice",
		SkillContent: skillTemplate,
		ExtraFiles:   files,
		Usage:        "install",
	}, args)
}

func collectTopicFiles() ([]install.InstallFile, error) {
	var files []install.InstallFile
	err := fs.WalkDir(topicsFS, topicsDir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		data, err := topicsFS.ReadFile(p)
		if err != nil {
			return fmt.Errorf("read embedded %s: %w", p, err)
		}
		files = append(files, install.InstallFile{
			Path:    p,
			Content: data,
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}
