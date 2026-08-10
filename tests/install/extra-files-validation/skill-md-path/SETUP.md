# Scenario

**Feature**: ExtraFiles path `"SKILL.md"` is rejected

```
HandleInstall(ExtraFiles Path="SKILL.md") -> cannot replace SKILL.md
```

## Preconditions
- An extra file has path `"SKILL.md"` (invalid — an extra file must not replace the main SKILL.md).

## Steps
1. Call `HandleInstall` with an extra file whose `Path` is `"SKILL.md"`.
2. The target directory argument is "test-skill".

```go
import "github.com/xhd2015/skills/install"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.ExtraFiles = []install.InstallFile{
		{Path: "SKILL.md", Content: []byte("content")},
	}
	req.Args = []string{"test-skill"}
	return nil
}
```
