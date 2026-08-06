// Command runupdate executes HandleInstall / HandleUpdate / HandleUpdateMany
// inside an isolated process so the parent doctest harness stays Parallel-safe
// (no os.Chdir / t.Setenv / os.Stdout reassignment in suite Run).
//
// Usage:
//
//	runupdate -req request.json
//
// WorkDir is created by the parent; this process chdirs into it, optionally
// rewrites HOME, runs PreInstalls (stdout discarded), applies mutations, then
// runs update with product stdout on this process's stdout.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/xhd2015/skills/install"
)

type fileMutate struct {
	RelPath string `json:"relPath"`
	Content string `json:"content"`
}

type preInstall struct {
	Opts install.InstallOptions `json:"opts"`
	Args []string               `json:"args"`
}

type request struct {
	UseMany           bool                  `json:"useMany"`
	SingleOpts        install.InstallOptions `json:"singleOpts"`
	ManySkills        []install.UpdateSkill `json:"manySkills"`
	UpdateArgs        []string              `json:"updateArgs"`
	PreInstalls       []preInstall          `json:"preInstalls"`
	PostInstallMutate []fileMutate          `json:"postInstallMutate"`
	UseGlobalHome     bool                  `json:"useGlobalHome"`
	WorkDir           string                `json:"workDir"`
	HomeDir           string                `json:"homeDir"`
}

func main() {
	reqPath := flag.String("req", "", "path to request JSON")
	flag.Parse()
	if *reqPath == "" {
		fmt.Fprintln(os.Stderr, "runupdate: -req is required")
		os.Exit(2)
	}
	data, err := os.ReadFile(*reqPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "runupdate: read req: %v\n", err)
		os.Exit(2)
	}
	var req request
	if err := json.Unmarshal(data, &req); err != nil {
		fmt.Fprintf(os.Stderr, "runupdate: decode req: %v\n", err)
		os.Exit(2)
	}
	if req.WorkDir == "" {
		fmt.Fprintln(os.Stderr, "runupdate: workDir is required")
		os.Exit(2)
	}
	if err := os.Chdir(req.WorkDir); err != nil {
		fmt.Fprintf(os.Stderr, "runupdate: chdir: %v\n", err)
		os.Exit(2)
	}
	if req.UseGlobalHome {
		if req.HomeDir == "" {
			fmt.Fprintln(os.Stderr, "runupdate: homeDir required when useGlobalHome")
			os.Exit(2)
		}
		if err := os.Setenv("HOME", req.HomeDir); err != nil {
			fmt.Fprintf(os.Stderr, "runupdate: set HOME: %v\n", err)
			os.Exit(2)
		}
	}

	for i, m := range req.PostInstallMutate {
		if strings.HasPrefix(m.RelPath, "$HOME/") {
			home, err := os.UserHomeDir()
			if err != nil {
				fmt.Fprintf(os.Stderr, "runupdate: home: %v\n", err)
				os.Exit(2)
			}
			req.PostInstallMutate[i].RelPath = filepath.Join(home, strings.TrimPrefix(m.RelPath, "$HOME/"))
		}
	}

	// Pre-install must not pollute captured update stdout.
	oldStdout := os.Stdout
	devNull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "runupdate: open devnull: %v\n", err)
		os.Exit(2)
	}
	os.Stdout = devNull
	for _, step := range req.PreInstalls {
		if err := install.HandleInstall(step.Opts, step.Args); err != nil {
			os.Stdout = oldStdout
			_ = devNull.Close()
			fmt.Fprintf(os.Stderr, "runupdate: pre-install: %v\n", err)
			os.Exit(1)
		}
	}
	os.Stdout = oldStdout
	_ = devNull.Close()

	for _, m := range req.PostInstallMutate {
		if err := os.MkdirAll(filepath.Dir(m.RelPath), 0755); err != nil {
			fmt.Fprintf(os.Stderr, "runupdate: mkdir: %v\n", err)
			os.Exit(1)
		}
		if err := os.WriteFile(m.RelPath, []byte(m.Content), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "runupdate: write mutate: %v\n", err)
			os.Exit(1)
		}
	}

	var runErr error
	if req.UseMany {
		runErr = install.HandleUpdateMany(req.ManySkills, req.UpdateArgs)
	} else {
		runErr = install.HandleUpdate(req.SingleOpts, req.UpdateArgs)
	}
	if runErr != nil {
		// Product error: surface on stderr; exit 1. Parent records Error field.
		fmt.Fprintln(os.Stderr, runErr.Error())
		os.Exit(1)
	}
	// Ensure product stdout is fully flushed (fmt uses os.Stdout).
	_, _ = io.WriteString(os.Stdout, "")
}
