# Scenario

**Feature**: FormatHeaderWithDelimiters re-wraps frontmatter with `---`

```
# header-only show path
caller -> skillcmd.FormatHeaderWithDelimiters(content) -> ---\n...\n---\n
```

## Steps

1. Set FileOp format-header with demo skill content.

```go
func Setup(t *testing.T, req *Request) error {
	req.FileOp = FileOpFormatHeader
	req.Content = demoRootContent
	return nil
}
```
