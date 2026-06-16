## Preconditions
- An extra file has path `".."` (invalid — escapes the target directory).

## Steps
1. Call `HandleInstall` with an extra file whose `Path` is `".."`.
2. The target directory argument is "test-skill".

```go
import "github.com/xhd2015/skills/install"

func Setup(t *testing.T, req *Request) error {
	req.ExtraFiles = []install.InstallFile{
		{Path: "..", Content: []byte("content")},
	}
	req.Args = []string{"test-skill"}
	return nil
}
```
