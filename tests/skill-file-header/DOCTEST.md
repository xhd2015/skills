# Skill File Header Parsing Tests

Doc-style tests for `skill_file.GetHeader`, `skill_file.ParseHeader`, and
`Entries.Get` — pure functions that extract and parse YAML frontmatter from
SKILL.md content.

# DSN (Domain Specific Notion)

Participants:

- A SKILL.md file stores optional YAML frontmatter between `---` delimiters
  followed by a Markdown body.
- `GetHeader` reads raw file content and returns the inner YAML text without
  delimiter lines, or an error when the opening `---\n` prefix is absent.
- `ParseHeader` turns YAML header text into ordered `Entry` name/value pairs,
  including folded block scalars such as `description: >-`.
- `Entries.Get` looks up a header field by exact key name and returns an empty
  string when the key is missing.

## Decision Tree

```text
skill-file-header
├── input-validity
│   ├── valid-header
│   └── missing-header
└── header-format
    └── folded-description
```

## Test Index

- `input-validity/valid-header`: valid SKILL.md returns inner YAML and parsed entries including `name=git-fetch`.
- `input-validity/missing-header`: content without `---\n` prefix makes `GetHeader` return an error.
- `header-format/folded-description`: folded `description: >-` block parses to a single normalized value.

## How to Run

```sh
doctest vet ./tests/skill-file-header
doctest test -v ./tests/skill-file-header
```

## Version

0.0.2

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
	"github.com/xhd2015/skills/skill_file"
)

type Request struct {
	Content string
}

type Response struct {
	Header     string
	HeaderErr  error
	Entries    skill_file.Entries
	ParseErr   error
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	_ = d
	header, headerErr := skill_file.GetHeader(req.Content)
	resp := &Response{
		Header:    header,
		HeaderErr: headerErr,
	}
	if headerErr == nil {
		entries, parseErr := skill_file.ParseHeader(header)
		resp.Entries = entries
		resp.ParseErr = parseErr
	}
	return resp, nil
}
```