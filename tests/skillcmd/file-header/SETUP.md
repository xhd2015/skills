# Scenario

**Feature**: skillcmd file-header APIs extract and reformat YAML frontmatter

```
# pure functions over SKILL.md content
caller -> GetHeader / FormatHeaderWithDelimiters -> YAML text | error
```

## Preconditions

- Mode is file-header; sample content set per leaf.

## Steps

1. Set `req.Mode = ModeFileHeader`.
2. Leaves set Content and FileOp.

```go
func Setup(t *testing.T, req *Request) error {
	req.Mode = ModeFileHeader
	return nil
}
```
