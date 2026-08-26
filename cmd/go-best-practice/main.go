// Command go-best-practice is a redirect stub.
//
// The full recipe CLI moved to https://github.com/xhd2015/go-best-practice
// (module github.com/xhd2015/go-best-practice). This package path remains
// installable for old bookmarks and always prints a move message on stderr,
// then exits 1.
package main

import (
	"fmt"
	"os"
)

const redirectMessage = `go-best-practice has moved and is maintained at:

  https://github.com/xhd2015/go-best-practice

  Module:  github.com/xhd2015/go-best-practice
  Install: go install github.com/xhd2015/go-best-practice/cmd/go-best-practice@latest

This skills package path only redirects; it no longer serves recipes or skill content.
`

func main() {
	fmt.Fprint(os.Stderr, redirectMessage)
	os.Exit(1)
}
