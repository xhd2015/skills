## Preconditions
- An extra file has path `"SKILL.md"` (invalid — an extra file must not replace the main SKILL.md).

## Steps
1. Call `HandleInstall` with an extra file whose `Path` is `"SKILL.md"`.
2. The target directory argument is "test-skill".

```go
import "github.com/xhd2015/skills/install"

func Setup(t *testing.T, req *Request) error {
	req.ExtraFiles = []install.InstallFile{
		{Path: "SKILL.md", Content: []byte("content")},
	}
	req.Args = []string{"test-skill"}
	return nil
}
```
