# Scenario

**Feature**: GetHeader returns inner YAML without delimiter lines

```
# extract frontmatter
caller -> skillcmd.GetHeader(content) -> name/description YAML
```

## Steps

1. Set FileOp get-header and sample content with delimiters.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.FileOp = FileOpGetHeader
	req.Content = "---\nname: git-fetch\ndescription: clone repos\n---\n\n# Body\n"
	return nil
}
```
