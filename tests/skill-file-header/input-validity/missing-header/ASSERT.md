## Expected

- `GetHeader` returns a non-nil error.
- `Header` is empty.
- `Entries` is nil or empty because parsing is skipped.

## Side Effects

- None.

## Errors

- `GetHeader` must fail with a descriptive error about missing or malformed header.

## Exit Code

- Not applicable.

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if resp.HeaderErr == nil {
		t.Fatal("GetHeader error = nil, want error for missing header")
	}
	if resp.Header != "" {
		t.Fatalf("Header = %q, want empty when GetHeader fails", resp.Header)
	}
	if len(resp.Entries) != 0 {
		t.Fatalf("Entries = %#v, want empty when GetHeader fails", resp.Entries)
	}
}
```