package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/xhd2015/dot-pkgs/go-pkgs/govet"
)

func handleVet(args []string) error {
	cfg := govet.Config{FileMaxLines: 500}
	var dirs []string
	var jsonOutput bool

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "--json":
			jsonOutput = true
		case arg == "--file-max-lines":
			if i+1 >= len(args) {
				return fmt.Errorf("--file-max-lines requires a value")
			}
			i++
			n, err := strconv.Atoi(args[i])
			if err != nil || n < 0 {
				return fmt.Errorf("--file-max-lines must be a non-negative integer, got %q", args[i])
			}
			cfg.FileMaxLines = n
		case strings.HasPrefix(arg, "--file-max-lines="):
			n, err := strconv.Atoi(strings.TrimPrefix(arg, "--file-max-lines="))
			if err != nil || n < 0 {
				return fmt.Errorf("--file-max-lines must be a non-negative integer, got %q", strings.TrimPrefix(arg, "--file-max-lines="))
			}
			cfg.FileMaxLines = n
		case arg == "--exclude":
			if i+1 >= len(args) {
				return fmt.Errorf("--exclude requires a value")
			}
			i++
			cfg.Excludes = append(cfg.Excludes, args[i])
		case strings.HasPrefix(arg, "--exclude="):
			cfg.Excludes = append(cfg.Excludes, strings.TrimPrefix(arg, "--exclude="))
		case strings.HasPrefix(arg, "-"):
			return fmt.Errorf("unknown flag: %s", arg)
		default:
			dirs = append(dirs, arg)
		}
	}

	if len(dirs) == 0 {
		dirs = []string{"."}
	}

	workDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("get working directory: %w", err)
	}

	resolved, err := govet.ResolvePatterns(workDir, dirs)
	if err != nil {
		return fmt.Errorf("resolve patterns: %w", err)
	}

	violations, err := govet.Run(cfg, resolved)
	if err != nil {
		return err
	}

	if jsonOutput {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(violations)
	}

	for _, v := range violations {
		fmt.Printf("%s:%d:%d: [%s] %s\n", v.File, v.Line, v.Col, v.Checker, v.Message)
		if v.Hint != "" {
			fmt.Printf("  → %s\n", v.Hint)
		}
	}
	return nil
}
