# skill-cli — building skill CLIs with `skill install` and `skill show`

Recipe for adding `skill install` and `skill show` sub-commands
to an existing CLI so users can install the skill into their
editor (e.g. Cursor) and view the skill definition.

## Directory structure

```
my-skill/
├── SKILL.md           # YAML frontmatter + markdown content
├── main.go            # CLI entry point
└── topics/            # embedded recipe files
    ├── foo.md
    └── bar/
        └── baz.md
```

`SKILL.md` contains a `name` and `description` in YAML front-matter,
followed by the skill's markdown body (topic index, usage examples,
etc.).

```yaml
---
name: my-skill
description: >-
  A skill that does something useful.
---

# My Skill

...
```

## Embedding resources

Use `//go:embed` to bundle `SKILL.md` and the `topics/` directory
into the binary:

```go
import (
    "embed"
    "io/fs"
)

//go:embed SKILL.md
var skillTemplate string

//go:embed topics
var topicsFS embed.FS
```

## `skill show` sub-command

Prints the embedded `SKILL.md` content to stdout. The handler is
straightforward — just print the embedded string:

```go
func handleSkill(args []string) error {
    if len(args) == 0 || args[0] != "show" {
        return fmt.Errorf("unknown skill sub-command: expected `skill show`")
    }
    fmt.Print(skillTemplate)
    return nil
}
```

## `skill install` sub-command

Uses `github.com/xhd2015/skills/install` to copy `SKILL.md` and
all topic files into a target directory (e.g. `.cursor/skills/my-skill`).

First, collect the topic files from the embedded filesystem:

```go
type installFile struct {
    Path    string
    Content []byte
}

func collectTopicFiles() ([]installFile, error) {
    var files []installFile
    err := fs.WalkDir(topicsFS, "topics", func(p string, d fs.DirEntry, err error) error {
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
        files = append(files, installFile{Path: p, Content: data})
        return nil
    })
    return files, err
}
```

Then wire up the install handler, delegating to the `install` package:

```go
import "github.com/xhd2015/skills/install"

func handleInstall(args []string) error {
    files, err := collectTopicFiles()
    if err != nil {
        return err
    }
    return install.HandleInstall(install.InstallOptions{
        SkillDirName: "my-skill",
        SkillContent: skillTemplate,
        ExtraFiles:   files,
        Usage:        "install",
    }, args)
}
```

The `install` package handles `--cursor` (installs into
`.cursor/skills/`) and custom `<dir>` targets.

## Putting it together

A minimal `main.go` that recognises `skill show`, `skill install`,
and topic lookup:

```go
package main

import (
    "embed"
    "fmt"
    "io/fs"
    "os"
    "path"
    "sort"
    "strings"

    "github.com/xhd2015/skills/install"
)

//go:embed SKILL.md
var skillTemplate string

//go:embed topics
var topicsFS embed.FS

func main() {
    if err := handle(os.Args[1:]); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}

func handle(args []string) error {
    if len(args) == 0 {
        fmt.Print(help)
        return printTopicIndex()
    }
    switch args[0] {
    case "install":
        return handleInstall(args[1:])
    case "skill":
        return handleSkill(args[1:])
    }
    content, ok, err := readTopic(args[0])
    if err != nil {
        return err
    }
    if !ok {
        return fmt.Errorf("unknown command or topic: %s", args[0])
    }
    fmt.Print(content)
    return nil
}

func handleSkill(args []string) error {
    if len(args) == 0 || args[0] != "show" {
        return fmt.Errorf("unknown skill sub-command: expected `skill show`")
    }
    fmt.Print(skillTemplate)
    return nil
}

func handleInstall(args []string) error {
    files, err := collectTopicFiles()
    if err != nil {
        return err
    }
    return install.HandleInstall(install.InstallOptions{
        SkillDirName: "my-skill",
        SkillContent: skillTemplate,
        ExtraFiles:   files,
        Usage:        "install",
    }, args)
}

func collectTopicFiles() ([]install.InstallFile, error) {
    var files []install.InstallFile
    err := fs.WalkDir(topicsFS, "topics", func(p string, d fs.DirEntry, err error) error {
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
        files = append(files, install.InstallFile{Path: p, Content: data})
        return nil
    })
    return files, err
}

func readTopic(topicPath string) (string, bool, error) {
    topicPath = strings.Trim(topicPath, "/")
    if topicPath == "" {
        return "", false, nil
    }
    segments := strings.Split(topicPath, "/")
    for _, s := range segments {
        if s == "" || s == "." || s == ".." {
            return "", false, fmt.Errorf("invalid topic path segment: %q", s)
        }
    }
    embedPath := path.Join(append([]string{"topics"}, segments...)...) + ".md"
    data, err := topicsFS.ReadFile(embedPath)
    if err != nil {
        if os.IsNotExist(err) {
            return "", false, nil
        }
        return "", false, fmt.Errorf("read topic %s: %w", topicPath, err)
    }
    return string(data), true, nil
}

func printTopicIndex() error {
    var topics []string
    fs.WalkDir(topicsFS, "topics", func(p string, d fs.DirEntry, err error) error {
        if err != nil {
            return err
        }
        if d.IsDir() || !strings.HasSuffix(p, ".md") {
            return nil
        }
        rel := strings.TrimPrefix(strings.TrimSuffix(p, ".md"), "topics/")
        topics = append(topics, rel)
        return nil
    })
    sort.Strings(topics)
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

const help = `
Usage: my-skill <command> [ARGS]
       my-skill <topic>[/<sub-topic>[/...]]

Commands:
  install [<dir>]    Install SKILL.md + topics to a directory (or use --cursor)
  skill show         Show the content of SKILL.md
  <topic-path>       Print the detailed content for a topic or sub-topic
`
```

## See also

- `flags-parsing/subcommand` — sub-command dispatcher patterns for
  the CLI dispatch in this recipe.
